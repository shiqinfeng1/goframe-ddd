package dockercmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type ComposeController struct {
}

func New(ctx context.Context) (*ComposeController, error) {
	return &ComposeController{}, nil
}

// docker images 所有镜像信息
func (ctl *ComposeController) LoadImage(ctx context.Context, imageFile string) error {

	return nil
}

// docker images 所有镜像信息
func (ctl *ComposeController) Images(ctx context.Context) ([]string, error) {
	reoptags := make([]string, 0)

	return reoptags, nil
}

// docker compose images 当前运行容器的镜像信息
func (ctl *ComposeController) ComposeImages(ctx context.Context) ([]string, error) {
	reoptags := make([]string, 0)
	return reoptags, nil
}
func (ctl *ComposeController) ComposeUp(ctx context.Context, version string) error {
	r, err := g.Client().Post(ctx, "http://host.docker.internal:31083/image/upgrade/"+version)
	if err != nil {
		return err
	}
	defer r.Close()
	g.Log().Infof(ctx, "docker upgrade image version:%v", r.ReadAllString())
	return nil
}
