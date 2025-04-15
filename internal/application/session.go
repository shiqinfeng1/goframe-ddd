package application

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/pkg/stream/session"
)

func (app *Application) GetClientIds(ctx context.Context) ([]string, error) {
	nodeIds, err := session.GetNodeList(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, nil
	}
	return nodeIds, nil
}
