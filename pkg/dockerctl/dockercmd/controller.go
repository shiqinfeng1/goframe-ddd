package dockercmd

import (
	"context"
	"log"

	"github.com/gogf/gf/v2/net/gclient"
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
	r, err := gclient.New().Post(ctx, "http://host.docker.internal:31083/image/upgrade/"+version)
	if err != nil {
		return err
	}
	defer r.Close()
	log.Println("docker upgrade image version:", r.ReadAllString())
	return nil
}
