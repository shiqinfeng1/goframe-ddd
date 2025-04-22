package ops

import (
	"context"
	"strings"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/ops/v1"
)

func (c *ControllerV1) ImageList(ctx context.Context, req *v1.ImageListReq) (res *v1.ImageListRes, err error) {
	res = &v1.ImageListRes{}
	images, err := c.dockerOps.Images(ctx)
	if err != nil {
		return nil, err
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
	res = &v1.ContainerImageRes{}
	images, err := c.dockerOps.ComposeImages(ctx)
	if err != nil {
		return nil, err
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
