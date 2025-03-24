package application

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/migration"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/pkg/stream"
)

type Application struct {
	fileTransfer FileTransferService
}

// New 一个DDD的应用层，支持CQRS，
func New(ctx context.Context) *Application {
	// 实例化一个文件管理的数据仓库
	repo := adapters.NewFilemgrRepo(migration.NewEntClient(ctx))
	// 实例化一个文件传输服务
	stm := stream.NewStream()
	maxTasks := g.Cfg().MustGet(ctx, "filemgr.maxTasks").Int()
	ftSrv := filemgr.NewFileTransferService(maxTasks, stm, repo)
	stm.Startup(ctx, ftSrv.StreamRecvHandler)

	return &Application{
		fileTransfer: ftSrv,
	}
}
