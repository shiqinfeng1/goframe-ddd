package ops

import (
	"context"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/ops/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) UpgradeApp(ctx context.Context, req *v1.UpgradeAppReq) (res *v1.UpgradeAppRes, err error) {
	res = &v1.UpgradeAppRes{}
	err = c.app.UpgradeApp(ctx, &application.UpgradeAppInput{
		AppName: req.AppName,
	})
	if err != nil {
		err = errors.ErrUpgradeAppFail(err)
	}
	return
}
func (c *ControllerV1) UpgradeImage(ctx context.Context, req *v1.UpgradeImageReq) (res *v1.UpgradeImageRes, err error) {
	res = &v1.UpgradeImageRes{}
	err = c.app.UpgradeImage(ctx, &application.UpgradeImageInput{
		Version: req.Version,
	})
	if err != nil {
		err = errors.ErrUpgradeImageFail(err)
	}
	return
}
