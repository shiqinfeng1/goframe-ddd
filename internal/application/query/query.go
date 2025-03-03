package query

import (
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
)

type Handler struct {
	repo filemgr.Repository
}

func NewHandler(
	repo filemgr.Repository,
) *Handler {
	return &Handler{
		repo: repo,
	}
}
