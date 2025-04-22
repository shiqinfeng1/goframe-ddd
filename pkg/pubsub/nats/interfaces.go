package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
	pubsub "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

//go:generate mockgen -destination=mock_client.go -package=nats -source=./interfaces.go ConnIntf,ConnMgr,SubMgr

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

// ConnMgr represents the main Client connection.
type ConnMgr interface {
	Connect(ctx context.Context, opts ...nats.Option)
	Close(ctx context.Context)
	Publish(ctx context.Context, subject string, message []byte) error
	JsPublish(ctx context.Context, subject string, message []byte) error
	GetJetStream() (jetstream.JetStream, error)
	isConnected() bool
	Health() *health.Health
}

// SubMgr represents the main Subscription Manager.
type JsSubMgr interface {
	New(ctx context.Context, stream streamIntf, identity []string, consumeType SubType) error
	Delete(ctx context.Context, identity []string) error
	Subscribe(ctx context.Context, identity []string, handler pubsub.JsSubscribeFunc) error
	Close(ctx context.Context) error
}
type SubMgr interface {
	New(ctx context.Context, conn *nats.Conn, topicName string, consumeType SubType) error
	Delete(ctx context.Context, topicName string) error
	Subscribe(ctx context.Context, topicName string, handler pubsub.SubscribeFunc) error
	Close(ctx context.Context) error
}
