package application

import (
	"context"
	"sync"

	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/service"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/infrastructure/repositories"
)

type Application struct {
	pointDataSet service.PointDataSetService
}

var app *Application
var once sync.Once

// New 一个DDD的应用层
func App(ctx context.Context) *Application {
	once.Do(func() {
		// 点位数据集服务
		repoPm := repositories.NewPointmgrRepo()
		pdsSrv := service.NewPointDataSetService(ctx, repoPm)

		app = &Application{
			pointDataSet: pdsSrv,
		}
	})
	return app
}
