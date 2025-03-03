package adapters

import (
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
)

type filemgrRepo struct{}

// NewTrainingRepo .
func NewFilemgrRepo() filemgr.Repository {
	return &filemgrRepo{}
}
