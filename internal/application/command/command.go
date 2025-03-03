package command

import (
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
)

type Handler struct {
	repo         filemgr.Repository
	fileTransfer FileTransferService
	stream       StreamService
}

func NewHandler(
	repo filemgr.Repository,
	fileTransfer FileTransferService,
	stream StreamService,
) *Handler {
	return &Handler{
		repo:         repo,
		fileTransfer: fileTransfer,
		stream:       stream,
	}
}
