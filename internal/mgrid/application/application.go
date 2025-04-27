package application

import (
	"context"

	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/service"
)

type Service struct {
	pointDataSet service.PointDataSetService
}

var WireProviderSet = wire.NewSet(New)

// New 一个DDD的应用层
func New(ctx context.Context, pdsSrv service.PointDataSetService) *Service {
	return &Service{
		pointDataSet: pdsSrv,
	}
}
