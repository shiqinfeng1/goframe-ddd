package dockercmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type DockerController struct {
}

func New(ctx context.Context) (*DockerController, error) {
	return &DockerController{}, nil
}

// docker images 所有镜像信息
func (ctl *DockerController) LoadImage(ctx context.Context, imageFile string) error {

	return nil
}

// docker images 所有镜像信息
func (ctl *DockerController) Images(ctx context.Context) ([]string, error) {
	reoptags := make([]string, 0)

	return reoptags, nil
}

// docker compose images 当前运行容器的镜像信息
func (ctl *DockerController) ComposeImages(ctx context.Context) ([]string, error) {
	reoptags := make([]string, 0)
	return reoptags, nil
}
func (ctl *DockerController) ComposeUp(ctx context.Context, version string) error {
	r, err := g.Client().Post(ctx, "http://host.docker.internal:31083/image/upgrade/"+version)
	if err != nil {
		return err
	}
	defer r.Close()
	g.Log().Infof(ctx, "docker upgrade image version:%v", r.ReadAllString())
	return nil
}
