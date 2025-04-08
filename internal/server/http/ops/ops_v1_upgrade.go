package ops

import (
	"context"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/ops/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
)

func (c *ControllerV1) Upgrade(ctx context.Context, req *v1.UpgradeReq) (res *v1.UpgradeRes, err error) {
	res = &v1.UpgradeRes{}
	err = c.app.UpgradeApp(ctx, &application.UpgradeAppInput{
		AppName: req.AppName,
	})
	return
}
