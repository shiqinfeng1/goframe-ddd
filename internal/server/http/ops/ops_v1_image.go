package ops

import (
	"context"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/ops/v1"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) ImageList(ctx context.Context, req *v1.ImageListReq) (res *v1.ImageListRes, err error) {
	res = &v1.ImageListRes{}
	out, err := c.app.Images(ctx)
	if err != nil {
		return res, errors.ErrQueryImageFail(err)
	}
	res.Images = out.Images
	return
}
func (c *ControllerV1) ContainerImage(ctx context.Context, req *v1.ContainerImageReq) (res *v1.ContainerImageRes, err error) {
	res = &v1.ContainerImageRes{}
	out, err := c.app.ComposeImages(ctx)
	if err != nil {
		return res, errors.ErrQueryImageFail(err)
	}
	res.Images = out.Images
	return
}
