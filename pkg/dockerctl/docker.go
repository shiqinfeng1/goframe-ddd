package dockerctl

import (
	"context"
)

type DockerOps interface {
	LoadImage(ctx context.Context, imageFile string) error
	Images(ctx context.Context) ([]string, error)
	ComposeImages(ctx context.Context) ([]string, error)
	ComposeUp(ctx context.Context, version string) error
}

type HandlerResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (hr *HandlerResponse) GetStrings() []string {
	if v, ok := hr.Data.([]string); ok {
		return v
	}
	return []string{}
}
