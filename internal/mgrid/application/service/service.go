package service

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type PointDataSetService interface {
	HandleTopic1(ctx context.Context, msg *nats.Msg) error
	HandleTopic2(ctx context.Context, msg *jetstream.Msg) error
}
