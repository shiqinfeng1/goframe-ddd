package repositories

import (
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/repository"
)

type pointmgrRepo struct {
}

var WireProviderSet = wire.NewSet(NewPointmgrRepo)

// NewTrainingRepo .
func NewPointmgrRepo() repository.Repository {
	return &pointmgrRepo{}
}
