package repositories

import (
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/repository"
)

type pointmgrRepo struct {
}

// NewTrainingRepo .
func NewPointmgrRepo() repository.Repository {
	return &pointmgrRepo{}
}
