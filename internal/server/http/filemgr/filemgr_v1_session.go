package filemgr

import (
	"context"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/filemgr/v1"
)

func (c *ControllerV1) SessionList(ctx context.Context, req *v1.SessionListReq) (res *v1.SessionListRes, err error) {
	out, err := c.app.Queries.GetClientIds(ctx)
	res = &v1.SessionListRes{
		ClientIds: out,
	}
	return
}
