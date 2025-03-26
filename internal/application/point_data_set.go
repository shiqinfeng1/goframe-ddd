package application

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

func (h *Application) HandleTopic1(ctx context.Context, msg *pubsub.Message) error {
	h.pointDataSet.HandleTopic1(ctx, msg)
	return nil
}
func (h *Application) HandleTopic2(ctx context.Context, msg *pubsub.Message) error {
	h.pointDataSet.HandleTopic2(ctx, msg)
	return nil
}
