package application

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func (h *Application) HandleTopic1(ctx context.Context, msg *nats.Msg) error {
	err := h.pointDataSet.HandleTopic1(ctx, msg)
	return err
}

func (h *Application) HandleTopic2(ctx context.Context, msg *jetstream.Msg) error {
	err := h.pointDataSet.HandleTopic2(ctx, msg)
	return err
}
