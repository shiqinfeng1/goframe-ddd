// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package filemgr

import (
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shiqinfeng1/goframe-ddd/api/http/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
)

type ControllerV1 struct {
	app *application.Application
}

func NewV1() filemgr.IFilemgrV1 {
	ctx := gctx.New()
	return &ControllerV1{
		app: application.App(ctx),
	}
}
