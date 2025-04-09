package composectl

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/api/types/image"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

var (
	errInvalidProject = fmt.Errorf("invalid docker compose project")
)

type ComposeController struct {
	project   *types.Project
	service   api.Service
	dockerCli *command.DockerCli
	logger    *logger
	running   atomic.Bool
}

func newProject(composePath string) (*types.Project, error) {
	p, err := ProjectFromConfig(composePath)
	if err != nil {
		return nil, err
	}
	if p.Name == "" {
		return nil, errInvalidProject
	}
	cfgs := projectPortConfigs(p)

	// 如果有端口配置冲突(同一端口被多个service映射)
	if hrojectPortConfigs(cfgs) {
		conflicts := portConflicts(cfgs) // 获取冲突的端口列表
		resolved, err := resolvePortConflicts(conflicts)
		if err != nil {
			return nil, err
		}
		applyPortMapping(p, resolved)
	}
	return p, nil
}
func New(ctx context.Context, composePath string) (*ComposeController, error) {

	project, err := newProject(composePath)
	if err != nil {
		return nil, err
	}

	c := &ComposeController{
		project: project,
	}

	logger, err := newLogConsumer(ctx)
	if err != nil {
		return nil, err
	}
	c.logger = logger

	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}
	err = dockerCli.Initialize(flags.NewClientOptions())
	if err != nil {
		return nil, err
	}
	c.dockerCli = dockerCli
	c.service = compose.NewComposeService(dockerCli)

	go func() {
		err := c.service.Logs(gctx.New(), project.Name, c.logger, api.LogOptions{
			Services:   nil,
			Tail:       "",
			Since:      "",
			Until:      "",
			Follow:     true,
			Timestamps: false,
		})
		if err != nil {
			g.Log().Errorf(ctx, "docker logs exit: %v", err)
			return
		}
		g.Log().Infof(ctx, "docker logs exit ok")
	}()

	c.running.Store(true)
	return c, nil
}

// docker images 所有镜像信息
func (ctl *ComposeController) LoadImage(ctx context.Context, imageFile string) error {
	file, err := gfile.Open(imageFile)
	if err != nil {
		return err
	}

	response, err := ctl.dockerCli.Client().ImageLoad(ctx, file)
	if err != nil {
		return err
	}
	return response.Body.Close()
}

// docker images 所有镜像信息
func (ctl *ComposeController) Images(ctx context.Context) ([]image.Summary, error) {
	images, err := ctl.dockerCli.Client().ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
}

// docker compose images 当前运行容器的镜像信息
func (ctl *ComposeController) ComposeImages(ctx context.Context) ([]api.ImageSummary, error) {
	images, err := ctl.service.Images(ctx, ctl.project.Name, api.ImagesOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
}
func (ctl *ComposeController) ComposeUp(ctx context.Context) error {
	err := ctl.service.Up(ctx, ctl.project, api.UpOptions{
		Start: api.StartOptions{
			Project: ctl.project,
		},
	})
	if err != nil {
		return err
	}

	ctl.running.Store(true)
	return nil
}

func (ctl *ComposeController) Down(ctx context.Context) error {
	if ctl.project == nil {
		return errInvalidProject
	}
	if !ctl.running.Load() {
		return nil
	}
	ctl.running.Store(false)

	err := ctl.service.Down(ctx, ctl.project.Name, api.DownOptions{
		Project:       ctl.project,
		RemoveOrphans: true,
	})
	if err != nil {
		return err
	}
	return nil
}
