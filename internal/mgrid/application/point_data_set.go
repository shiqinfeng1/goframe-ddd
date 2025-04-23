package application

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func (app *Application) HandleTopic1(ctx context.Context, msg *nats.Msg) error {
	err := app.pointDataSet.HandleTopic1(ctx, msg)
	return err
}

func (app *Application) HandleTopic2(ctx context.Context, msg *jetstream.Msg) error {
	err := app.pointDataSet.HandleTopic2(ctx, msg)
	return err
}
