package application

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/migration"
	"github.com/shiqinfeng1/goframe-ddd/internal/application/command"
	"github.com/shiqinfeng1/goframe-ddd/internal/application/query"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/pkg/stream"
)

type Application struct {
	Commands *command.Handler
	Queries  *query.Handler
}

// New 一个DDD的应用层，支持CQRS，
func New(ctx context.Context) *Application {
	// 实例化一个文件管理的数据仓库
	repo := adapters.NewFilemgrRepo(migration.NewEntClient(ctx))

	stm := stream.NewStream()

	// 实例化一个文件传输服务
	maxTasks := g.Cfg().MustGet(ctx, "filemgr.maxTasks").Int()
	ftSrv := filemgr.NewFileTransferService(maxTasks, stm, repo)
	stm.Startup(ctx, ftSrv.StreamRecvHandler)

	return &Application{
		Commands: command.NewHandler(ftSrv),
		Queries:  query.NewHandler(repo),
	}
}
