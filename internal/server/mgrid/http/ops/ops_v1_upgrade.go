package ops

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/ops/v1"
)

func (c *ControllerV1) UpgradeApp(ctx context.Context, req *v1.UpgradeAppReq) (res *v1.UpgradeAppRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) UpgradeImage(ctx context.Context, req *v1.UpgradeImageReq) (res *v1.UpgradeImageRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
