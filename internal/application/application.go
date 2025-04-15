package application

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/migration"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/pointmgr"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl"

	// "github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl/dockersock"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl/dockercmd"
	"github.com/shiqinfeng1/goframe-ddd/pkg/stream"
)

type Application struct {
	fileTransfer FileTransferService
	pointDataSet PointDataSetService
	dockerOps    dockerctl.DockerOps
}

var app *Application
var once sync.Once

// New 一个DDD的应用层
func App(ctx context.Context) *Application {
	once.Do(func() {
		// 实例化一个文件管理的数据仓库
		var ftSrv *filemgr.FileTransferMgr
		if g.Cfg().MustGet(ctx, "filemgr.enable").Bool() {
			repoFm := adapters.NewFilemgrRepo(migration.NewEntClient(ctx))
			stm := stream.NewStream() // 实例化一个文件传输服务
			maxTasks := g.Cfg().MustGet(ctx, "filemgr.maxTasks").Int()
			ftSrv = filemgr.NewFileTransferService(maxTasks, stm, repoFm)
			stm.Startup(ctx, ftSrv.StreamRecvHandler)
		}
		// 点位数据集服务
		repoPm := adapters.NewPointmgrRepo(migration.NewEntClient(ctx))
		pdsSrv := pointmgr.NewPointDataSetService(repoPm)

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
