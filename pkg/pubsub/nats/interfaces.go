package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
	pubsub "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

//go:generate mockgen -destination=mock_client.go -package=nats -source=./interfaces.go ConnIntf,ConnectionManagerIntf,SubscriptionManagerIntf

// ConnIntf represents the main Client connection.
type ConnIntf interface {
	Status() nats.Status
	Close()
	Conn() *nats.Conn
	NewJetStream() (jetstream.JetStream, error)
}

// Connector represents the main Client connection.
type Connector interface {
	Connect(string, ...nats.Option) (ConnIntf, error)
}

// JetStreamCreator represents the main Client jStream Client.
type JetStreamCreator interface {
	New(conn ConnIntf) (jetstream.JetStream, error)
}

// ConnectionManagerIntf represents the main Client connection.
type ConnectionManagerIntf interface {
	Connect(ctx context.Context)
	Close(ctx context.Context)
	Publish(ctx context.Context, subject string, message []byte) error
	JsPublish(ctx context.Context, subject string, message []byte) error
	GetJetStream() (jetstream.JetStream, error)
	isConnected() bool
	Health() *health.Health
}

// SubscriptionManagerIntf represents the main Subscription Manager.
type JsSubscriptionManagerIntf interface {
	NewSubscriber(stream streamIntf, streamName, consumerName, topicName string, consumeType SubType) *jsSubscriber
	DeleteSubscriber(ctx context.Context, streamName, consumerName, topicName string) error
	Subscribe(ctx context.Context, streamName, consumerName, topicName string, handler pubsub.SubscribeFunc) error
	Close(ctx context.Context) error
}
type SubscriptionManagerIntf interface {
	NewSubscriber(conn *nats.Conn, topicName string, consumeType SubType) *subscriber
	DeleteSubscriber(ctx context.Context, topicName string) error
	Subscribe(ctx context.Context, topicName string, handler pubsub.SubscribeFunc) error
	Close(ctx context.Context) error
}
