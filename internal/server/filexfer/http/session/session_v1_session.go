package session

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/shiqinfeng1/goframe-ddd/api/filexfer/http/session/v1"
)

func (c *ControllerV1) NodeList(ctx context.Context, req *v1.NodeListReq) (res *v1.NodeListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
