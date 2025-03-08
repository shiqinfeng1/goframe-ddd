package command

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
)

// 文件传输服务：文件读写，分块，发送，接收，任务队列
type FileTransferService interface {
	AddTask(ctx context.Context, taskId string, name string, nodeId string, path []string)
	PauseTask(ctx context.Context, taskId string)
	ResumeTask(ctx context.Context, taskId string)
	CancelTask(ctx context.Context, taskId string)
	GetTaskStatus(ctx context.Context, taskId string) filemgr.Status
}
