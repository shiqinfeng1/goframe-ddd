package application

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/application/service"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/pointmgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories/migration"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl"

	// "github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl/dockersock"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl/dockercmd"
)

type Application struct {
	fileTransfer service.FileTransferService
	pointDataSet service.PointDataSetService
	dockerOps    dockerctl.DockerOps
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

		// 实例化一个dockeecompose 控制器
		dockerOps, err := dockercmd.New(ctx)
		if err != nil {
			g.Log().Fatal(ctx, err)
		}
		app = &Application{
			fileTransfer: ftSrv,
			pointDataSet: pdsSrv,
			dockerOps:    dockerOps,
		}
	})
	return app
}
