package command

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/pkg/stream"
	"github.com/xtaci/smux"
)

// 文件传输服务：文件读写，分块，发送，接收，任务队列
type FileTransferService interface {
	AddTask(ctx context.Context, id string, name string, path []string)
	PauseTask(ctx context.Context, id string)
	ResumeTask(ctx context.Context, id string)
	GetTaskStatus(ctx context.Context, id string) filemgr.Status
}

// 数据流管理
type StreamService interface {
	OpenStreamByClient(ctx context.Context, handler stream.SendStreamHandleFunc) error
	OpenStreamByServer(ctx context.Context, session *smux.Session, handler stream.SendStreamHandleFunc) error
}
