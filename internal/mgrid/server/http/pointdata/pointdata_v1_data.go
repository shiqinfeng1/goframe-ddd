package pointdata

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/pointdata/v1"
)

func (c *ControllerV1) GetPointData(ctx context.Context, req *v1.GetPointDataReq) (res *v1.GetPointDataRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
