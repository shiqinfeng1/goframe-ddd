package application

import (
	"context"
	"sync"

	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/application/service"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/infrastructure/repositories"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/infrastructure/repositories/migration"
)

type Application struct {
	fileTransfer service.FilexferService
}

var app *Application
var once sync.Once

// New 一个DDD的应用层
func App(ctx context.Context) *Application {
	once.Do(func() {
		// 文件传输服务
		repoFm := repositories.NewFilemgrRepo(migration.NewEntClient(ctx))
		ftSrv := filemgr.NewFileTransferService(ctx, repoFm)

		ftSrv.Start(ctx)

		app = &Application{
			fileTransfer: ftSrv,
		}
	})
	return app
}
