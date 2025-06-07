package repositories

import (
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/repository"
)

type pointdataRepo struct {
}

// NewTrainingRepo .
func NewPointdataRepo() repository.PointdataRepository {
	return &pointdataRepo{}
}
