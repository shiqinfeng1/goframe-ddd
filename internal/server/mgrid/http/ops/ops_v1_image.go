package ops

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/ops/v1"
)

func (c *ControllerV1) ImageList(ctx context.Context, req *v1.ImageListReq) (res *v1.ImageListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) ContainerImage(ctx context.Context, req *v1.ContainerImageReq) (res *v1.ContainerImageRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
