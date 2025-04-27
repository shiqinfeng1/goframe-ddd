package application

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/pkg/session"
)

func (app *Service) GetClientIds(ctx context.Context) ([]string, error) {
	nodeIds, err := session.GetSessionNodeList(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, nil
	}
	return nodeIds, nil
}
