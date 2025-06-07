package dockercmd

import (
	"context"

	"github.com/gogf/gf/v2/net/gclient"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl"
)

type DockerController struct {
}

func New(ctx context.Context) (dockerctl.DockerOps, error) {
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
	// 访问宿主机上的watcher守护进程，进行版本升级
	r, err := gclient.New().Post(ctx, "http://host.docker.internal:31083/image/upgrade/"+version)
	if err != nil {
		return err
	}
	defer r.Close()
	return nil
}
