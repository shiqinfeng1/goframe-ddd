package filemgr

import (
	"context"
	"os"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gfile"
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
	repo          Repository
}

func (fs *fileSaver) Close() error {
	if fs == nil {
		return nil
	}
	if fs.inst != nil {
		fs.inst.Close()
		fs.inst = nil
	}

	return nil
}

func (fs *fileSaver) saveChunk(fc *fileChunk) error {
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
func (fs *fileSaver) monitorCancel(ctx context.Context, fileId string) {
	tickerSubCancelNotify := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tickerSubCancelNotify.C: // 监控是否有取消发送的通知
			fs.Close() // 发送方取消发送，退出和关闭filersaver
			g.Log().Infof(ctx, "cancel recv file sucess. fileId:%v", fileId)
			tickerSubCancelNotify.Stop()
			return
		case <-fs.timeoutTicker.C: // 强制超时时间5分钟关闭fileSaver，防止对端异常退出，filesaver资源无法正常释放，设置一个超时时间
			fs.Close()
			g.Log().Infof(ctx, "force stop recv file sucess")
			fs.timeoutTicker.Stop()
			return
		case <-ctx.Done():
			g.Log().Infof(ctx, "exit ferry monitor cancel ok")
			return
		}
	}
}

func (fs *fileSaver) SaveChunk(fc *fileChunk) {
	if err := fs.saveChunk(fc); err != nil {
		return
	}
	if fs.timeoutTicker != nil {
		fs.timeoutTicker.Reset(5 * time.Minute)
	}
}

// 一个文件对应一个fileSaver,
func NewFileSave(ctx context.Context, cache *gcache.Cache, taskId, fileId string, repo Repository) (*fileSaver, error) {
	// 检查缓存，如果存在，直接返回，如果不存在，新建
	val, err := cache.GetOrSetFuncLock(ctx, fileId, gcache.Func(func(ctx context.Context) (value interface{}, err error) {
		// 查询文件接收记录
		recvFile, err := repo.GetRecvTaskFile(ctx, taskId, fileId)
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
		}

		go fs.monitorCancel(ctx, fileId)
		return fs, nil
	}), 0)
	if err != nil {
		return nil, gerror.Wrapf(err, "get fileSaver from cache fail: taskId=%v fileId=%v", taskId, fileId)
	}
	var fs *fileSaver
	err = val.Scan(&fs)
	if err != nil || fs == nil {
		return nil, gerror.Newf("get file saver fail:%v", err)
	}
	return fs, nil
}
