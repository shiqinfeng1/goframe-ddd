package filemgr

import (
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/pkg/utils"
	"github.com/xtaci/smux"
	"golang.org/x/net/context"
)

type (
	RecvStreamHandleFunc func(context.Context, *smux.Session, io.ReadWriter) error
	SendStreamHandleFunc func(context.Context, *smux.Stream) error
)

// 数据流管理
type StreamIntf interface {
	SendByClient(ctx context.Context, handler SendStreamHandleFunc) error
	SendByServer(ctx context.Context, session *smux.Session, handler SendStreamHandleFunc) error
}

// 任务状态枚举
type Status struct {
	val  int
	desc string
}

var (
	StatusUndefined  = Status{val: 0, desc: "未定义"}
	StatusWaiting    = Status{val: 1, desc: "等待发送"}
	StatusSending    = Status{val: 2, desc: "正在发送"}
	StatusPaused     = Status{val: 3, desc: "已暂停"}
	StatusCancel     = Status{val: 4, desc: "已取消"}
	StatusFailed     = Status{val: 5, desc: "发送失败"}
	StatusSuccessful = Status{val: 6, desc: "发送成功"}
)

type (
	postSendFunc func(bool)
	postFunc     func()
)

// FileSendTask 表示文件发送任务
type TransferTask struct {
	taskName                 string
	taskId                   string
	nodeId                   string
	paths                    []string
	status                   Status
	sendFileChan             chan postSendFunc
	pauseChan                chan postFunc
	cancelChan               chan postFunc
	stream                   StreamIntf
	chunkOffsets, chunkSizes []int64
	notifyStop               atomic.Bool
}

func NewTransferTask(ctx context.Context, id, name, nodeId string, paths []string, stream StreamIntf) *TransferTask {
	task := &TransferTask{
		taskId:       id,
		paths:        paths,
		taskName:     name,
		nodeId:       nodeId,
		status:       StatusWaiting,
		sendFileChan: make(chan postSendFunc, 4),
		pauseChan:    make(chan postFunc, 4),
		cancelChan:   make(chan postFunc, 4),
		stream:       stream,
	}
	go task.worker(ctx)
	return task
}

func (t *TransferTask) openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, gerror.Wrapf(err, "open file fail:%v", filePath)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, gerror.Wrapf(err, "get file stat fail:%v", filePath)
	}
	t.chunkOffsets, t.chunkSizes, err = utils.SplitFile(info.Size())
	if err != nil {
		file.Close()
		return nil, gerror.Wrapf(err, "get file chunk fail:%v", filePath)
	}
	return file, nil
}

func (t *TransferTask) sendChunk(ctx context.Context, file *os.File, stm *smux.Stream) error {
	for i, chunkSize := range t.chunkSizes {
		if t.notifyStop.Load() {
			t.notifyStop.Store(false)
			return nil
		}
		// todo 构造header
		header := []byte{}
		// 分块获取文件数据,串行发送
		section := io.NewSectionReader(file, t.chunkOffsets[i], chunkSize)
		start := time.Now()
		n, err := stm.Write(header)
		if err != nil {
			return gerror.Wrap(err, "header write stream fail")
		}
		if n != len(header) {
			return gerror.Newf("write stream fail:n(%v) != len(header)(%v)", n, len(header))
		}
		written, err := io.CopyN(stm, section, section.Size()) // s.Write([]byte(msg))
		if err != nil {
			return gerror.Wrap(err, "write stream fail")
		}
		end := time.Since(start)
		g.Log().Debugf(ctx, "stresm write %v bytes ok.  stream.id=%v elapsed:%v write-speed:%v MB/s", written, stm.ID(), end, float64(written)/1024/1024/end.Seconds())

		recvd := make([]byte, 1024)
		m, err := stm.Read(recvd)
		if err != nil {
			return gerror.Wrap(err, "read stream fail")
		}
		end2 := time.Since(start)
		g.Log().Debugf(ctx, "client read resp %v bytes.  stream.id=%v elapsed:%v roundtrip-speed:%v MB/s", m, stm.ID(), end2, float64(written)/1024/1024/end2.Seconds())
		// todo 文件和分块存储到数据库
	}
	return nil
}

func (t *TransferTask) worker(ctx context.Context) {
	for {
		select {
		case postHandle := <-t.cancelChan:
			t.notifyStop.Store(true)
			postHandle()
		case postHandle := <-t.pauseChan:
			t.notifyStop.Store(true)
			postHandle()
		case postHandle := <-t.sendFileChan:
			doSend := func(ctx context.Context, stm *smux.Stream) error {
				// 遍历所有待发送的文件
				for _, filePath := range t.paths {
					if t.notifyStop.Load() {
						t.notifyStop.Store(false)
						break
					}
					// 打开要发送的文件，如果失败，记录到本地，但不同步到对端
					file, err := t.openFile(filePath)
					if err != nil {
						postHandle(false)
						return err
					}
					// defer file.Close()
					// todo 文件和分块存储到数据库
					// todo 发送文件信息给对端

					if err := t.sendChunk(ctx, file, stm); err != nil {
						postHandle(false)
						return err
					}
				}
				postHandle(true)
				return nil
			}

			if t.nodeId != "" { // 服务端需指定nodeid
				sess, err := Session().GetSession(ctx, t.nodeId)
				if err != nil {
					g.Log().Errorf(ctx, "send file fail:%v", err)
					postHandle(false)
					return
				}
				if err := t.stream.SendByServer(ctx, sess, doSend); err != nil {
					g.Log().Errorf(ctx, "send file by server fail:%v", err)
					postHandle(false)
					return
				}
			} else {
				if err := t.stream.SendByClient(ctx, doSend); err != nil {
					g.Log().Errorf(ctx, "send file by client fail:%v", err)
					postHandle(false)
					return
				}
			}
		}
	}
}
