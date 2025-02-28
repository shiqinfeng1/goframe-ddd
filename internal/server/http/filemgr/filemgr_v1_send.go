package filemgr

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/shiqinfeng1/goframe-ddd/api/http/filemgr/v1"
)

func (c *ControllerV1) SendFile(ctx context.Context, req *v1.SendFileReq) (res *v1.SendFileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
