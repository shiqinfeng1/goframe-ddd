package filemgr

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/filemgr/v1"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) NodeList(ctx context.Context, req *v1.NodeListReq) (res *v1.NodeListRes, err error) {
	if !g.Cfg().MustGet(ctx, "filemgr.enable").Bool() {
		return res, errors.ErrFileMgrDisable
	}
	out, err := c.app.GetClientIds(ctx)
	res = &v1.NodeListRes{
		NodeIds: out,
	}
	return
}
