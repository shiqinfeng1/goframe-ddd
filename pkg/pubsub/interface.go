package pubsub

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type JsSubscribeFunc func(ctx context.Context, msg *jetstream.Msg) error
type SubscribeFunc func(ctx context.Context, msg *nats.Msg) error

type Committer interface {
	Commit(ctx context.Context)
}
