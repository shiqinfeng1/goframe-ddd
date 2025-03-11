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
	offset     int64
	data       []byte
	chunkIndex uint32
	md5        string
}

type fileSaver struct {
	fileId        string
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
func (fs *fileSaver) monitorEvent(ctx context.Context) {
	for {
		select {
		case status := <-fs.eventNotify: // 监控事件通知(取消、暂停)
			// 2. 删除文件
			if status == StatusCancel {
				g.Log().Infof(ctx, "TODO: 任务已取消，需删除文件:%v", fs.path)
			}
			// 更新recvfile状态 为取消或暂停
			if err := fs.repo.UpdateRecvStatus(ctx, fs.fileId, status); err != nil {
				g.Log().Errorf(ctx, "update recvfile fileId:%v status fail:%v", fs.fileId, err)
				return
			}
			// 删除实例
			if err := removeFileSaver(ctx, fs.fileId); err != nil {
				g.Log().Errorf(ctx, "remove filesaver fileId:%v fail:%v", fs.fileId, err)
				return
			}
			g.Log().Infof(ctx, "update fileId:%v to status:%v ok", fs.fileId, status)
			return
		case <-fs.timeoutTicker.C: // 强制超时时间5分钟关闭fileSaver，防止对端异常退出，filesaver资源无法正常释放，设置一个超时时间
			if err := fs.repo.UpdateRecvStatus(ctx, fs.fileId, StatusFailed); err != nil {
				g.Log().Errorf(ctx, "update recvfile fileId:%v status to fail:%v", fs.fileId, err)
				return
			}
			if err := removeFileSaver(ctx, fs.fileId); err != nil {
				g.Log().Errorf(ctx, "force stop recv file fileId:%v fail:%v", fs.fileId, err)
				return
			}
			g.Log().Errorf(ctx, "force stop recv file success. fileId:%v", fs.fileId)
			return
		}
	}
}

func (fs *fileSaver) SaveChunk(ctx context.Context, fc *fileChunk) error {
	// 每个文件块超时5分钟， 如果超时未收到下一个块，主动释放资源
	if fs.timeoutTicker != nil {
		fs.timeoutTicker.Reset(5 * time.Minute)
	}
	if err := fs.saveChunkData(fc); err != nil {
		removeFileSaver(ctx, fs.fileId)
		return err
	}
	newrf, err := fs.repo.UpdateRecvChunk(ctx, &RecvChunk{
		FileID:      fs.fileId,
		ChunkIndex:  int(fc.chunkIndex),
		ChunkOffset: fc.offset,
		ChunkSize:   len(fc.data),
	})
	if err != nil {
		removeFileSaver(ctx, fs.fileId)
		return gerror.Wrapf(err, "save recvfile fail")
	}
	// g.Log().Debugf(ctx, "recv file chunk[%v/%v]:%v", newrf.ChunkNumRecved, newrf.ChunkNumTotal, fs.path)
	if NewStatus(uint32(newrf.Status)) == StatusSuccessful {
		fs.inst.Close()
		fs.inst = nil
		if err := gfile.Rename(fs.path+".downloading", fs.path); err != nil {
			return gerror.Wrapf(err, "rename recvfile fail:%v", fs.path)
		}
		g.Log().Infof(ctx, "recv file ok:%v", fs.path)
		removeFileSaver(ctx, fs.fileId)
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
		if recvFile.FilePathSave == "" {
			return nil, gerror.Newf("not found recvfile:%v", fileId)
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
			fileId:        fileId,
			path:          recvFile.FilePathSave, // 使用实际保存的路径，如果重名，文件名可能被重命名
			timeoutTicker: time.NewTicker(5 * time.Minute),
			inst:          file,
			repo:          repo,
			eventNotify:   make(chan Status),
		}

		go fs.monitorEvent(ctx)
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
	// 释放filesaver资源
	fs.Close()
	return nil
}
