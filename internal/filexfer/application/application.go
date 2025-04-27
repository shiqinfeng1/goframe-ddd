package application

import (
	"context"

	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/application/service"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/domain/filemgr"
)

type Service struct {
	fileTransfer service.FilexferService
}

var WireProviderSet = wire.NewSet(New)

// New 一个DDD的应用层
func New(ctx context.Context, ftSrv *filemgr.FileTransferMgr) *Service {

	ftSrv.Start(ctx)

	return &Service{
		fileTransfer: ftSrv,
	}

}
