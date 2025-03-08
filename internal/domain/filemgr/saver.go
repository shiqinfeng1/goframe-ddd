package filemgr

import (
	"context"
	"os"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/shiqinfeng1/goframe-ddd/pkg/cache"
)

type fileChunk struct {
	taskId     string
	fileId     string
	offset     int64
	data       []byte
	chunkIndex uint32
	md5        string
}

type fileSaver struct {
	inst          *os.File
	path          string
	timeoutTicker *time.Ticker
	eventNotify   chan Status
	repo          Repository
}

func (fs *fileSaver) Close() {
	if fs == nil {
		return
	}
	if fs.inst != nil {
		fs.inst.Close()
		fs.inst = nil
	}
	if fs.timeoutTicker != nil {
		fs.timeoutTicker.Stop()
		fs.timeoutTicker = nil
	}
}

func (fs *fileSaver) saveChunkData(fc *fileChunk) error {
	n, err := fs.inst.WriteAt(fc.data, fc.offset)
	if err != nil {
		return gerror.Wrapf(err, "Error writing to file")
	} else if n < len(fc.data) {
		return gerror.Newf("Only wrote %d bytes; expected %d\n", n, len(fc.data))
	}
	// return gerror.Newf("Only wrote %d bytes; expected %d\n", n, len(fc.data))
	return nil
}

// 监听对端取消发送的事件：发送方取消之后，需要通知接收方释放资源
func (fs *fileSaver) monitorEvent(ctx context.Context, fileId string) {
	for {
		select {
		case status := <-fs.eventNotify: // 监控是否有取消发送的通知
			// 1. 关闭 filesaver 实例
			fs.Close() // 发送方取消发送，退出和关闭filersaver
			// 2. 删除文件
			if status == StatusCancel {
				// todo 删除文件
			}
			if status == StatusPaused {
				g.Log().Infof(ctx, "pause recv file success. fileId:%v", fileId)
			}
			// todo 更新recvfile状态 为取消或暂停
			// 删除实例
			removeFileSaver(ctx, fileId)
			g.Log().Infof(ctx, "cancel recv file success. fileId:%v", fileId)
			return
		case <-fs.timeoutTicker.C: // 强制超时时间5分钟关闭fileSaver，防止对端异常退出，filesaver资源无法正常释放，设置一个超时时间
			fs.Close()
			// todo 更新recvfile状态 为 中断
			removeFileSaver(ctx, fileId)
			g.Log().Infof(ctx, "force stop recv file success")
			return
		}
	}
}

func (fs *fileSaver) SaveChunk(ctx context.Context, fc *fileChunk) error {
	if fs.timeoutTicker != nil {
		fs.timeoutTicker.Reset(5 * time.Minute)
	}
	if err := fs.saveChunkData(fc); err != nil {
		return err
	}
	finished, err := fs.repo.UpdateRecvChunk(ctx, &RecvChunk{
		FileID:      fc.fileId,
		ChunkIndex:  int(fc.chunkIndex),
		ChunkOffset: fc.offset,
		ChunkSize:   len(fc.data),
	})
	if err != nil {
		return gerror.Wrapf(err, "save recvfile fail")
	}

	if finished {
		fs.Close()
		removeFileSaver(ctx, fc.fileId)
	}
	return nil
}

func (fs *fileSaver) EventNotify(status int) {
	fs.eventNotify <- NewStatus(uint32(status))
}

// 一个文件对应一个fileSaver,
func getFileSaver(ctx context.Context, fileId string, repo Repository) (*fileSaver, error) {
	// 检查缓存，如果存在，直接返回，如果不存在，新建
	val, err := cache.Memory().GetOrSetFuncLock(ctx, fileId, gcache.Func(func(ctx context.Context) (value interface{}, err error) {
		// 查询文件接收记录
		recvFile, err := repo.GetRecvFile(ctx, fileId)
		if err != nil {
			return nil, err
		}
		// 指定的文件已存在， 但是对应的downloading文件不存在，那么不需要新建对应的downloading文件
		if gfile.IsFile(recvFile.FilePathSave) {
			if !gfile.IsFile(recvFile.FilePathSave + ".downloading") {
				g.Log().Warningf(ctx, "%v is already exist! change to <.downloading> file", recvFile.FilePathSave)
				if err := gfile.Rename(recvFile.FilePathSave, recvFile.FilePathSave+".downloading"); err != nil {
					return nil, err
				}
			}
		}
		// 新建或打开文件
		file, err := os.OpenFile(recvFile.FilePathSave+".downloading", os.O_CREATE|os.O_RDWR, 0o666)
		if err != nil {
			return nil, err
		}

		// 设置文件大小
		if err := file.Truncate(recvFile.FileSize); err != nil {
			return nil, err
		}
		// 块数据保存的管理结构
		fs := &fileSaver{
			path:          recvFile.FilePathSave, // 使用实际保存的路径，如果重名，文件名可能被重命名
			timeoutTicker: time.NewTicker(5 * time.Minute),
			inst:          file,
			repo:          repo,
			eventNotify:   make(chan Status),
		}

		go fs.monitorEvent(ctx, fileId)
		return fs, nil
	}), 0)
	if err != nil {
		return nil, gerror.Wrapf(err, "get fileSaver from cache fail: fileId=%v", fileId)
	}
	var fs *fileSaver
	err = val.Scan(&fs)
	if err != nil || fs == nil {
		return nil, gerror.Newf("get file saver fail:%v", err)
	}
	return fs, nil
}

func mustGetFileSaver(ctx context.Context, fileId string) (*fileSaver, error) {
	if exist, _ := cache.Memory().Contains(ctx, fileId); !exist {
		return nil, nil
	}
	val, err := cache.Memory().Get(ctx, fileId)
	if err != nil {
		return nil, gerror.Wrapf(err, "get fileSaver from cache fail: fileId=%v", fileId)
	}
	var fs *fileSaver
	err = val.Scan(&fs)
	if err != nil || fs == nil {
		return nil, gerror.Newf("must get file saver fail:%v", err)
	}
	return fs, nil
}

func removeFileSaver(ctx context.Context, fileId string) error {
	val, err := cache.Memory().Remove(ctx, fileId)
	if err != nil {
		return gerror.Wrapf(err, "get fileSaver from cache fail: fileId=%v", fileId)
	}
	var fs *fileSaver
	err = val.Scan(&fs)
	if err != nil || fs == nil {
		return gerror.Newf("must get file saver fail:%v", err)
	}
	fs.Close()
	return nil
}
