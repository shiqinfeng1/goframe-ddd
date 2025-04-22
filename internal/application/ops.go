package application

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/shiqinfeng1/goframe-ddd/internal/application/dto"
)

func (app *Application) UpgradeApp(ctx context.Context, in *dto.UpgradeAppInput) error {
	go func() {
		gproc.ShellExec(gctx.New(), `supervisorctl restart `+in.AppName)
	}()
	g.Log().Debugf(ctx, "upgrade app '%v' ...", in.AppName)

	return nil
}
func (app *Application) UpgradeImage(ctx context.Context, in *dto.UpgradeImageInput) error {
	go func() {
		nctx := gctx.New()
		if err := app.dockerOps.ComposeUp(nctx, in.Version); err != nil {
			g.Log().Debugf(nctx, "exec docker compose up fail':%v", err)
		}
	}()
	g.Log().Debugf(ctx, "upgrade image from '' to '%v' ...", in.Version)

	return nil
}

func (app *Application) ComposeImages(ctx context.Context) (*dto.ComposeImagesOutput, error) {
	images, err := app.dockerOps.ComposeImages(ctx)
	if err != nil {
		return nil, err
	}

	out := &dto.ComposeImagesOutput{
		Images: make([]dto.ImageSummary, 0),
	}
	for _, repotag := range images {
		repotags := strings.Split(repotag, ":")
		out.Images = append(out.Images, dto.ImageSummary{
			Name: repotags[0],
			Tag:  repotags[1],
		})
	}

	return out, nil
}
func (app *Application) Images(ctx context.Context) (*dto.ImagesOutput, error) {
	images, err := app.dockerOps.Images(ctx)
	if err != nil {
		return nil, err
	}

	out := &dto.ImagesOutput{
		Images: make([]dto.ImageSummary, 0),
	}
	for _, repotag := range images {
		repotags := strings.Split(repotag, ":")
		out.Images = append(out.Images, dto.ImageSummary{
			Name: repotags[0],
			Tag:  repotags[1],
		})
	}

	return out, nil
}
