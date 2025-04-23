package auth

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/auth/v1"
)

func (c *ControllerV1) DeviceAuth(ctx context.Context, req *v1.DeviceAuthReq) (res *v1.DeviceAuthRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
