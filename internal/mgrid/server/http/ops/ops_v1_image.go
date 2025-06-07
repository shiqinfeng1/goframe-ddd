package ops

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	v1 "github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/ops/v1"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) ImageList(ctx context.Context, req *v1.ImageListReq) (res *v1.ImageListRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	res = &v1.ImageListRes{}
	images, err := c.dockerOps.Images(ctx)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrGetImageListFail(lang)
	}

	out := &v1.ImagesOutput{
		Images: make([]v1.ImageSummary, 0),
	}
	for _, repotag := range images {
		repotags := strings.Split(repotag, ":")
		out.Images = append(out.Images, v1.ImageSummary{
			Name: repotags[0],
			Tag:  repotags[1],
		})
	}
	res.Images = out.Images
	return
}
func (c *ControllerV1) ContainerImage(ctx context.Context, req *v1.ContainerImageReq) (res *v1.ContainerImageRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	res = &v1.ContainerImageRes{}
	images, err := c.dockerOps.ComposeImages(ctx)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrGetImageListFail(lang)
	}

	out := &v1.ComposeImagesOutput{
		Images: make([]v1.ImageSummary, 0),
	}
	for _, repotag := range images {
		repotags := strings.Split(repotag, ":")
		out.Images = append(out.Images, v1.ImageSummary{
			Name: repotags[0],
			Tag:  repotags[1],
		})
	}

	res.Images = out.Images
	return
}
