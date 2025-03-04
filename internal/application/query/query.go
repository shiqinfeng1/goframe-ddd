package query

import (
	"context"

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

func (h *Handler) GetClientIds(ctx context.Context) ([]string, error) {
	clientsIds, err := filemgr.Session().GetSessionList(ctx)
	if err != nil {
		return nil, nil
	}
	return clientsIds, nil
}
