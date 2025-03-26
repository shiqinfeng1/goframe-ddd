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
	Connect(ctx context.Context) error
	CreateTopic(ctx context.Context, subjects []string) error
	DeleteTopic(ctx context.Context, topicName string) error

	Close(ctx context.Context) error
}

type Committer interface {
	Commit(ctx context.Context)
}
