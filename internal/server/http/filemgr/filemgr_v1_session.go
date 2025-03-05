package filemgr

import (
	"context"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/filemgr/v1"
)

func (c *ControllerV1) SessionList(ctx context.Context, req *v1.NodeListReq) (res *v1.NodeListRes, err error) {
	out, err := c.app.Queries.GetClientIds(ctx)
	res = &v1.NodeListRes{
		NodeIds: out,
	}
	return
}
