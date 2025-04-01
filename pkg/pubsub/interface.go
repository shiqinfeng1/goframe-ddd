package pubsub

import (
	"context"
)

type SubscribeFunc func(ctx context.Context, msg *Message) error

type Publisher interface {
	Publish(ctx context.Context, topic string, message []byte) error
	JsPublish(ctx context.Context, subject string, message []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler SubscribeFunc) error
	JsSubscribe(ctx context.Context, topic string, handler SubscribeFunc) error
}

type Committer interface {
	Commit(ctx context.Context)
}
