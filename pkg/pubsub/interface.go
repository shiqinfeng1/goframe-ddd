package pubsub

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, message []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (*Message, error)
}

type Client interface {
	Publisher
	Subscriber

	CreateTopic(ctx context.Context, name string) error
	DeleteTopic(ctx context.Context, name string) error

	Close() error
}

type Committer interface {
	Commit(ctx context.Context)
}
