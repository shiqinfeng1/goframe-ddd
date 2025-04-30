package pubsub

import (
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
)

func NewSubOrConsume(logger server.Logger, app application.Service) *ControllerV1 {
	subMgr := NewV1(logger, app)
	return subMgr
}
