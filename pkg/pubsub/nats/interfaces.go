package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

//go:generate mockgen -destination=mock_client.go -package=nats -source=./interfaces.go Client,Subscription,ConnIntf,ConnectionManagerIntf,SubscriptionManagerIntf,StreamManagerInterface

// ConnIntf represents the main Client connection.
type ConnIntf interface {
	Status() nats.Status
	Close()
	NATSConn() *nats.Conn
	JetStream() (jetstream.JetStream, error)
}

// Connector represents the main Client connection.
type Connector interface {
	Connect(string, ...nats.Option) (ConnIntf, error)
}

// JetStreamCreator represents the main Client jStream Client.
type JetStreamCreator interface {
	New(conn ConnIntf) (jetstream.JetStream, error)
}

// JetStreamClient represents the main Client jStream Client.
type JetStreamClient interface {
	Publish(ctx context.Context, subject string, message []byte) error
	Subscribe(ctx context.Context, subject string, handler messageHandler) error
	Close(ctx context.Context) error
	DeleteStream(ctx context.Context, name string) error
	CreateStream(ctx context.Context, cfg StreamConfig) error
	CreateOrUpdateStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error)
}

// ConnectionManagerIntf represents the main Client connection.
type ConnectionManagerIntf interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context)
	Publish(ctx context.Context, subject string, message []byte) error
	jetStream() (jetstream.JetStream, error)
	isConnected() bool
	Health() health.Health
}

// SubscriptionManagerIntf represents the main Subscription Manager.
type SubscriptionManagerIntf interface {
	Subscribe(
		ctx context.Context,
		topic string,
		js jetstream.JetStream,
		cfg *Config) (*pubsub.Message, error)
	Close()
}

// StreamManagerInterface represents the main Stream Manager.
type StreamManagerInterface interface {
	CreateStream(ctx context.Context, cfg StreamConfig) error
	DeleteStream(ctx context.Context, name string) error
	CreateOrUpdateStream(ctx context.Context, cfg *jetstream.StreamConfig) (jetstream.Stream, error)
}
