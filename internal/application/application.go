package application

import (
	"context"
	"sync"

	"github.com/shiqinfeng1/goframe-ddd/internal/application/service"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/pointmgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories/migration"
	// "github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl/dockersock"
)

type Application struct {
	fileTransfer service.FileTransferService
	pointDataSet service.PointDataSetService
}

var app *Application
var once sync.Once

// New 一个DDD的应用层
func App(ctx context.Context) *Application {
	once.Do(func() {
		// 文件传输服务
		repoFm := repositories.NewFilemgrRepo(migration.NewEntClient(ctx))
		ftSrv := filemgr.NewFileTransferService(ctx, repoFm)
		// 点位数据集服务
		repoPm := repositories.NewPointmgrRepo(migration.NewEntClient(ctx))
		pdsSrv := pointmgr.NewPointDataSetService(ctx, repoPm)

		ftSrv.Start(ctx)

		app = &Application{
			fileTransfer: ftSrv,
			pointDataSet: pdsSrv,
		}
	})
	return app
}
