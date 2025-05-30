// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package ops

import (
	"github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/ops"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl"
)

type ControllerV1 struct {
	logger    server.Logger
	app       application.Service
	dockerOps dockerctl.DockerOps
}

func NewV1(logger server.Logger, app application.Service, dockerOps dockerctl.DockerOps) ops.IOpsV1 {
	return &ControllerV1{
		logger:    logger,
		app:       app,
		dockerOps: dockerOps,
	}
}
