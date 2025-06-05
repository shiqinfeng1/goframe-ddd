package application

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/internal/demo/application/service"
	"github.com/shiqinfeng1/goframe-ddd/internal/demo/domain/entity/filemgr"
)

type Service struct {
	fileTransfer service.FilexferService
}

// New 一个DDD的应用层
func New(ctx context.Context, ftSrv *filemgr.FileTransferMgr) *Service {

	ftSrv.Start(ctx)

	return &Service{
		fileTransfer: ftSrv,
	}

}
