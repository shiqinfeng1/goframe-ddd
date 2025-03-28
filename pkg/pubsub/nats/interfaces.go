package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
	pubsub "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

//go:generate mockgen -destination=mock_client.go -package=nats -source=./interfaces.go ConnIntf,ConnectionManagerIntf,SubscriptionManagerIntf,StreamManagerIntf

// ConnIntf represents the main Client connection.
type ConnIntf interface {
	Status() nats.Status
	Close()
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
	Connect(ctx context.Context) error
	Close(ctx context.Context)
	Publish(ctx context.Context, subject string, message []byte) error
	getJetStream() (jetstream.JetStream, error)
	isConnected() bool
	Health() *health.Health
}

// SubscriptionManagerIntf represents the main Subscription Manager.
type SubscriptionManagerIntf interface {
	Subscribe(
		ctx context.Context,
		topic string,
		js jetstream.JetStream,
		cfg *Config,
		handler pubsub.SubscribeFunc,
	) error
	Close()
}

// StreamManagerIntf represents the main Stream Manager.
type StreamManagerIntf interface {
	CreateStream(ctx context.Context, cfg StreamConfig) error
	DeleteStream(ctx context.Context, name string) error
}
