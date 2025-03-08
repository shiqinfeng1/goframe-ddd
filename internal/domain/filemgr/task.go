package filemgr

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/rs/xid"
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
	val  uint32
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

func NewStatus(v uint32) Status {
	switch v {
	case StatusWaiting.val:
		return StatusWaiting
	case StatusSending.val:
		return StatusSending
	case StatusPaused.val:
		return StatusPaused
	case StatusCancel.val:
		return StatusCancel
	case StatusFailed.val:
		return StatusFailed
	case StatusSuccessful.val:
		return StatusSuccessful
	}
	return StatusUndefined
}

type (
	postSendFunc func(bool)
	postFunc     func()
)

// FileSendTask 表示文件发送任务
type TransferTask struct {
	repo         Repository
	taskName     string
	taskId       string
	nodeId       string
	paths        []string
	status       Status
	sendFileChan chan postSendFunc
	pauseChan    chan postFunc
	cancelChan   chan postFunc
	stream       StreamIntf
	chunkOffsets []int64
	chunkSizes   []int
	notifyStatus atomic.Uint32
}

func NewTransferTask(ctx context.Context, id, name, nodeId string, paths []string, stream StreamIntf, repo Repository) *TransferTask {
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
		repo:         repo,
	}
	go task.worker(ctx)
	return task
}

func (t *TransferTask) String() string {
	return fmt.Sprintf("taskId:%v taskName:%v nodeId:%v paths:%+v", t.taskId, t.taskName, t.nodeId, t.paths)
}

func (t *TransferTask) updateStatusAndNotifyPeer(ctx context.Context, status Status, stm io.ReadWriter) error {
	if err := t.repo.UpdateSendStatus(ctx, t.taskId, status); err != nil {
		return gerror.Wrapf(err, "update send status fail")
	}
	if err := t.syncEventToPeer(ctx, status, stm); err != nil {
		return err
	}
	return nil
}

func (t *TransferTask) getFileAndChunks(ctx context.Context, filePath string, stream io.ReadWriter) (*SendFile, *os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, gerror.Wrapf(err, "open file fail:%v", filePath)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, nil, gerror.Wrapf(err, "get file stat fail:%v", filePath)
	}

	// 从数据库读取断点信息
	sendFile, err := t.repo.GetSendFile(ctx, t.taskId, filePath)
	if err != nil {
		file.Close()
		return nil, nil, gerror.Wrapf(err, "get file info from repo fail:%v", filePath)
	}
	// 切割文件块
	t.chunkOffsets, t.chunkSizes, err = utils.SplitFile(info.Size())
	if err != nil {
		file.Close()
		return nil, nil, gerror.Wrapf(err, "get file chunk fail:%v", filePath)
	}
	// 数据库如果没有记录，说明需要发送全部文件, 同步信息给对端并记录这些文件信息到数据库，
	// 注意： 同步成功才记录到数据库
	if sendFile == nil {
		sendFile = &SendFile{
			TaskID:         t.taskId,
			TaskName:       t.taskName,
			FilePath:       filePath,
			FileId:         xid.New().String(),
			FileSize:       info.Size(),
			ChunkNumTotal:  len(t.chunkOffsets),
			ChunkNumSended: 0,
			Status:         -1,
		}
		if err := t.syncFileInfoToPeer(ctx, sendFile, stream); err != nil {
			file.Close()
			return nil, nil, gerror.Wrapf(err, "sync file info fail:%+v", sendFile)
		}
		// 文件和分块更新到数据库
		recid, err := t.repo.SaveSendFile(ctx, sendFile)
		if err != nil {
			file.Close()
			return nil, nil, gerror.Wrapf(err, "save send file fail:%v", filePath)
		}
		sendFile.ID = recid
		return sendFile, file, nil
	}
	// 该文件已发送完成
	if sendFile.ChunkNumSended == sendFile.ChunkNumTotal && sendFile.ChunkNumSended == len(t.chunkOffsets) {
		file.Close()
		return sendFile, nil, nil
	}
	// 数据库已记录文件块信息
	if sendFile.ChunkNumTotal == len(t.chunkOffsets) && sendFile.ChunkNumSended < len(t.chunkOffsets) {
		t.chunkOffsets = t.chunkOffsets[sendFile.ChunkNumSended:]
		t.chunkSizes = t.chunkSizes[sendFile.ChunkNumSended:]
	}

	return sendFile, file, nil
}

func (t *TransferTask) syncFileInfoToPeer(ctx context.Context, sendFile *SendFile, stm io.ReadWriter) error {
	body, _ := json.Marshal(sendFile)
	fiBytes := fileInfoMsgToBytes(ctx, body)
	// 分块获取文件数据,串行发送
	_, err := stm.Write(fiBytes)
	if err != nil {
		return gerror.Wrap(err, "fileinfo write stream fail")
	}
	// 接收响应数据
	respBody, err := recvAck(ctx, stm, msgFileInfo)
	if err != nil {
		return err
	}
	filePath := gconv.String(respBody)
	if sendFile.FilePath != filePath {
		return gerror.Newf("sync fileinfo fail. not match filepath: exp:%v fact:%v", sendFile.FilePath, filePath)
	}
	return nil
}

func (t *TransferTask) syncEventToPeer(ctx context.Context, status Status, stm io.ReadWriter) error {
	body, _ := json.Marshal(&EventMsg{
		TaskId: t.taskId,
		Status: int(status.val),
	})
	fiBytes := fileEventMsgToBytes(ctx, body)
	// 分块获取文件数据,串行发送
	_, err := stm.Write(fiBytes)
	if err != nil {
		return gerror.Wrap(err, "fileinfo write stream fail")
	}
	// 接收响应数据
	respBody, err := recvAck(ctx, stm, msgFileEvent)
	if err != nil {
		return err
	}
	taskId := gconv.String(respBody)
	if t.taskId != taskId {
		return gerror.Newf("sync fileinfo fail. not match taskid: exp:%v fact:%v", t.taskId, taskId)
	}
	return nil
}

func (t *TransferTask) sendChunk(ctx context.Context, sendFile *SendFile, file *os.File, stm io.ReadWriter) (bool, error) {
	for i, chunkSize := range t.chunkSizes {
		status := NewStatus(t.notifyStatus.Load())
		if status == StatusCancel || status == StatusPaused {
			if err := t.updateStatusAndNotifyPeer(ctx, status, stm); err != nil {
				return true, err
			}
			return true, nil
		}
		body, _ := json.Marshal(&SendChunk{
			TaskID:      sendFile.TaskID,
			FileID:      sendFile.FileId,
			ChunkIndex:  sendFile.ChunkNumSended,
			ChunkOffset: t.chunkOffsets[i],
			ChunkSize:   chunkSize,
		})
		fcBytes := fileChunkMsgToBytes(ctx, body)
		// 分块获取文件数据,串行发送
		section := io.NewSectionReader(file, t.chunkOffsets[i], int64(chunkSize))
		// 再次检查读取的数据是否一致
		if int64(chunkSize) != section.Size() {
			return false, gerror.Newf("read chunk fail: exp(%v) fac(%v)", chunkSize, section.Size())
		}
		start := time.Now()
		n, err := stm.Write(fcBytes)
		if err != nil {
			return false, gerror.Wrap(err, "filechunk header write fail")
		}
		if n != len(fcBytes) {
			return false, gerror.Newf("write stream fail:n(%v) != len(header)(%v)", n, len(fcBytes))
		}
		written, err := io.CopyN(stm, section, section.Size()) // s.Write([]byte(msg))
		if err != nil {
			return false, gerror.Wrap(err, "filechunk data write fail")
		}
		end := time.Since(start)
		g.Log().Debugf(ctx, "stresm write %v bytes ok. elapsed:%v write-speed:%v MB/s",
			written, end, float64(written)/1024/1024/end.Seconds())

		respBody, err := recvAck(ctx, stm, msgFileChunk)
		if err != nil {
			return false, gerror.Wrap(err, "filechunk ack fail")
		}
		end2 := time.Since(start)
		g.Log().Debugf(ctx, "client read ack %v bytes. elapsed:%v roundtrip-speed:%v MB/s",
			len(respBody), end2, float64(written)/1024/1024/end2.Seconds())

		fileId := gconv.String(respBody)
		if sendFile.FileId != fileId {
			return false, gerror.Newf("send filechunk fail. not match fileId: exp:%v fact:%v", sendFile.FilePath, fileId)
		}
		// 收到确认后， 更新本地数据库
		if err := t.repo.UpdateSendChunk(ctx, &SendChunk{
			TaskID:      sendFile.TaskID,
			FileID:      sendFile.FileId,
			SendFileID:  sendFile.ID, // 关联sendfile的i主键d
			ChunkIndex:  i,
			ChunkOffset: t.chunkOffsets[i],
			ChunkSize:   chunkSize,
		}); err != nil {
			return false, gerror.Wrap(err, "filechunk save fail")
		}
	}
	return false, nil
}

func (t *TransferTask) worker(ctx context.Context) {
	finishChan := make(chan struct{})
	for {
		select {
		case <-finishChan:
			g.Log().Infof(ctx, "exit file-tranfer task by finish:%v", t)
		case postHandle := <-t.cancelChan:
			t.notifyStatus.Store(StatusCancel.val)
			postHandle()
			g.Log().Infof(ctx, "exit file-tranfer task by cancel:%v", t)
		case postHandle := <-t.pauseChan:
			t.notifyStatus.Store(StatusPaused.val)
			postHandle()
			g.Log().Infof(ctx, "exit file-tranfer task by pause:%v", t)
		case postHandle := <-t.sendFileChan:
			t.notifyStatus.Store(StatusSending.val)

			// 执行发送任务
			doSend := func(ctx context.Context, stm *smux.Stream) error {
				defer close(finishChan)
				// 遍历所有待发送的文件
				for _, filePath := range t.paths {
					status := NewStatus(t.notifyStatus.Load())
					if status == StatusCancel || status == StatusPaused {
						if err := t.updateStatusAndNotifyPeer(ctx, status, stm); err != nil {
							return err
						}
						return nil // 直接返回，不需要执行postHandle，因为cancel和pause已执行各自的postHandle
					}
					// 打开要发送的文件，如果失败，记录到本地，但不同步到对端
					sendFile, file, err := t.getFileAndChunks(ctx, filePath, stm)
					if err != nil {
						postHandle(false)
						return err
					}
					if file == nil {
						g.Log().Infof(ctx, "file is already sended:%v skip it!", filePath)
						continue
					}
					// 发送文件块
					interrupt, err := t.sendChunk(ctx, sendFile, file, stm)
					if err != nil {
						file.Close()
						postHandle(false)
						return err
					}
					if interrupt { // 在发送分块期间被取消或暂停
						file.Close()
						return nil // 直接返回，不需要执行postHandle，因为cancel和pause已执行各自的postHandle
					}
					file.Close()
				}
				// 所有文件发送完成
				postHandle(true)
				return nil
			}

			// 在协程中执行回调doSend
			// 在回调doSend中执行成功后调用postHandle
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
