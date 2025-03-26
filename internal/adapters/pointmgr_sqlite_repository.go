package adapters

import (
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent"
)

type pointmgrRepo struct {
	db *ent.Client
}

// NewTrainingRepo .
func NewPointmgrRepo(db *ent.Client) *pointmgrRepo {
	return &pointmgrRepo{
		db: db,
	}
}
