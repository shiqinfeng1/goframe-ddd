package query

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
)

// 文件传输服务：文件读写，分块，发送，接收，任务队列
type FileTransferService interface {
	GetMaxAndRunning(ctx context.Context) (running int, maxtasks int)
	GetNotCompletedTasks(ctx context.Context) ([]*filemgr.FileTransferTask, map[string][]*filemgr.SendFile, error)
	GetCompletedTasks(ctx context.Context) ([]*filemgr.FileTransferTask, map[string][]*filemgr.SendFile, error)
}
