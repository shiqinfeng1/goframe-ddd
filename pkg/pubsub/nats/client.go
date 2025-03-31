package nats

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
	pubsub "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

const defaultRetryTimeout = 10 * time.Second

// Client represents a Client for NATS jStream operations.
type Client struct {
	connManager      ConnectionManagerIntf
	jsManager        JsSubscriptionManagerIntf
	subManager       SubscriptionManagerIntf
	Config           *Config
	natsConnector    Connector
	jetStreamCreator JetStreamCreator
}

// New creates a new Client.
func New(cfg *Config) *Client {
	if cfg == nil {
		cfg = &Config{}
	}

	client := &Client{
		Config: cfg,
	}

	return client
}

// Connect establishes a connection to NATS and sets up jStream.
func (c *Client) Connect(ctx context.Context) error {
	if c.connManager != nil && c.connManager.Health().Status == health.StatusUp {
		g.Log().Warning(ctx, "NATS connection already established")
		return nil
	}

	if c.Config.Server == "" {
		return errServerNotProvided
	}
	if c.Config.ConsumerName == "" {
		return errConsumerNotProvided
	}

	c.connManager = newConnectionManager(c.Config, c.natsConnector, c.jetStreamCreator)
	c.jsManager = NewJsSubscriberManager()
	c.subManager = NewSubscriberManager()

	c.connManager.Connect(ctx)

	return nil
}

func (c *Client) JetStream() (jetstream.JetStream, error) {
	if c.connManager == nil || c.connManager.Health().Status == health.StatusDown {
		return nil, gerror.New("NATS connection not established")
	}
	js, err := c.connManager.GetJetStream()
	if err != nil {
		return nil, err
	}
	return js, nil
}
func (c *Client) Conn() (*nats.Conn, error) {
	js, err := c.JetStream()
	if err != nil {
		return nil, err
	}
	return js.Conn(), nil
}

// Publish publishes a message to a topic.
func (c *Client) Publish(ctx context.Context, subject string, message []byte) error {
	return c.connManager.Publish(ctx, subject, message)
}
func (c *Client) JsPublish(ctx context.Context, subject string, message []byte) error {
	return c.connManager.JsPublish(ctx, subject, message)
}

// Subscribe subscribes to a topic and returns a single message.
func (c *Client) JsSubscribe(ctx context.Context, streamName, consumerName, topicName string, handler pubsub.SubscribeFunc) error {
	if !c.connManager.isConnected() {
		time.Sleep(defaultRetryTimeout)
		return errClientNotConnected
	}

	if err := c.jsManager.Subscribe(ctx, streamName, consumerName, topicName, handler); err != nil {
		return err
	}
	return nil
}
func (c *Client) Subscribe(ctx context.Context, topicName string, handler pubsub.SubscribeFunc) error {
	if !c.connManager.isConnected() {
		time.Sleep(defaultRetryTimeout)
		return errClientNotConnected
	}

	if err := c.subManager.Subscribe(ctx, topicName, handler); err != nil {
		return err
	}
	return nil
}

// Close closes the Client.
func (c *Client) Close(ctx context.Context) error {
	if c.subManager != nil {
		c.subManager.Close(ctx)
		c.subManager = nil
	}
	if c.connManager != nil {
		c.connManager.Close(ctx)
		c.connManager = nil
	}
	g.Log().Infof(ctx, "nats client close ok")
	return nil
}
