package composectl

import (
	"context"
)

type DockerOps interface {
	LoadImage(ctx context.Context, imageFile string) error
	Images(ctx context.Context) ([]string, error)
	ComposeImages(ctx context.Context) ([]string, error)
	ComposeUp(ctx context.Context, version string) error
}
