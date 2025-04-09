package application

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
)

func (app *Application) UpgradeApp(ctx context.Context, in *UpgradeAppInput) error {
	go func() {
		gproc.ShellExec(gctx.New(), `supervisorctl restart `+in.AppName)
	}()
	g.Log().Debugf(ctx, "upgrade app '%v' ...", in.AppName)

	return nil
}
func (app *Application) UpgradeImage(ctx context.Context, in *UpgradeImageInput) error {
	if err := app.imageController.ComposeUp(ctx); err != nil {
		return err
	}
	g.Log().Debugf(ctx, "upgrade image from '%v' to ...", in.Version)

	return nil
}

func (app *Application) ComposeImages(ctx context.Context) (*ComposeImagesOutput, error) {
	images, err := app.imageController.ComposeImages(ctx)
	if err != nil {
		return nil, err
	}

	out := &ComposeImagesOutput{
		Images: make([]ImageSummary, 0),
	}
	for _, image := range images {
		out.Images = append(out.Images, ImageSummary{
			Name: image.Repository,
			Tag:  image.Tag,
		})
	}

	return out, nil
}
func (app *Application) Images(ctx context.Context) (*ImagesOutput, error) {
	images, err := app.imageController.Images(ctx)
	if err != nil {
		return nil, err
	}

	out := &ImagesOutput{
		Images: make([]ImageSummary, 0),
	}
	for _, image := range images {
		for _, repotag := range image.RepoTags {
			repotags := strings.Split(repotag, ":")
			out.Images = append(out.Images, ImageSummary{
				Name: repotags[0],
				Tag:  repotags[1],
			})
		}
	}

	return out, nil
}
