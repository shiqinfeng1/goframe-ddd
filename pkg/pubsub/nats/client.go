package nats

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

//go:generate mockgen -destination=mock_tracer.go -package=nats go.opentelemetry.io/otel/trace Tracer

const defaultRetryTimeout = 10 * time.Second

var errClientNotConnected = errors.New("nats client not connected")

// Client represents a Client for NATS jStream operations.
type Client struct {
	connManager      ConnectionManagerIntf
	subManager       SubscriptionManagerIntf
	subscriptions    map[string]context.CancelFunc
	subMutex         sync.Mutex
	streamManager    StreamManagerInterface
	Config           *Config
	natsConnector    Connector
	jetStreamCreator JetStreamCreator
}

type messageHandler func(context.Context, jetstream.Msg) error

// New creates a new Client.
func New(cfg *Config) *Client {
	if cfg == nil {
		cfg = &Config{}
	}

	client := &Client{
		Config:     cfg,
		subManager: newSubscriptionManager(batchSize),
	}

	return client
}

// Connect establishes a connection to NATS and sets up jStream.
func (c *Client) Connect(ctx context.Context) error {
	g.Log().Debugf(ctx, "connecting to NATS server at %v", c.Config.Server)

	if err := c.validateAndPrepare(ctx); err != nil {
		return err
	}

	connManager := newConnectionManager(ctx, c.Config, c.natsConnector, c.jetStreamCreator)
	if err := connManager.Connect(ctx); err != nil {
		g.Log().Errorf(ctx, "failed to connect to NATS server at %v: %v", c.Config.Server, err)
		return err
	}

	c.connManager = connManager

	js, err := c.connManager.jetStream()
	if err != nil {
		return err
	}

	c.streamManager = newStreamManager(js)
	c.subManager = newSubscriptionManager(batchSize)
	c.logSuccessfulConnection(ctx)

	return nil
}

func (c *Client) validateAndPrepare(ctx context.Context) error {
	if err := validateConfigs(c.Config); err != nil {
		g.Log().Errorf(ctx, "could not initialize NATS jStream: %v", err)

		return err
	}

	return nil
}

func (c *Client) logSuccessfulConnection(ctx context.Context) {
	g.Log().Debugf(ctx, "connected to NATS server '%s'", c.Config.Server)
}

// Publish publishes a message to a topic.
func (c *Client) Publish(ctx context.Context, subject string, message []byte) error {
	return c.connManager.Publish(ctx, subject, message)
}

// Subscribe subscribes to a topic and returns a single message.
func (c *Client) Subscribe(ctx context.Context, topic string) (*pubsub.Message, error) {
	for {
		if !c.connManager.isConnected() {
			time.Sleep(defaultRetryTimeout)

			return nil, errClientNotConnected
		}

		js, err := c.connManager.jetStream()
		if err == nil {
			return c.subManager.Subscribe(ctx, topic, js, c.Config)
		}

		g.Log().Debugf(ctx, "Waiting for NATS connection to be established for topic %s", topic)

		time.Sleep(defaultRetryTimeout)
	}
}

func (c *Client) generateConsumerName(subject string) string {
	// return fmt.Sprintf("%s_%s", c.Config.Consumer, strings.ReplaceAll(subject, ".", "_"))
	return c.Config.Consumer
}

func (c *Client) SubscribeWithHandler(ctx context.Context, subject string, handler messageHandler) error {
	c.subMutex.Lock()
	defer c.subMutex.Unlock()

	// Cancel any existing subscription for this subject
	c.cancelExistingSubscription(subject)

	js, err := c.connManager.jetStream()
	if err != nil {
		return err
	}

	consumerName := c.generateConsumerName(subject)

	cons, err := c.createOrUpdateConsumer(ctx, js, subject, consumerName)
	if err != nil {
		return err
	}

	// Create a new context for this subscription
	subCtx, cancel := context.WithCancel(ctx)
	c.subscriptions[subject] = cancel

	go func() {
		defer cancel() // Ensure the cancellation is handled properly
		c.processMessages(subCtx, cons, subject, handler)
	}()

	return nil
}

func (c *Client) cancelExistingSubscription(subject string) {
	if cancel, exists := c.subscriptions[subject]; exists {
		cancel()
		delete(c.subscriptions, subject)
	}
}

func (c *Client) createOrUpdateConsumer(
	ctx context.Context, js jetstream.JetStream, subject, consumerName string,
) (jetstream.Consumer, error) {
	cons, err := js.CreateOrUpdateConsumer(ctx, c.Config.Stream.Stream, jetstream.ConsumerConfig{
		Durable:       consumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: subject,
		MaxDeliver:    c.Config.Stream.MaxDeliver,
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		g.Log().Errorf(ctx, "failed to create or update consumer: %v", err)
		return nil, err
	}

	return cons, nil
}

func (c *Client) processMessages(ctx context.Context, cons jetstream.Consumer, subject string, handler messageHandler) {
	for ctx.Err() == nil {
		if err := c.fetchAndProcessMessages(ctx, cons, subject, handler); err != nil {
			g.Log().Errorf(ctx, "Error in message processing loop for subject %s: %v", subject, err)
		}
	}
}

func (c *Client) fetchAndProcessMessages(ctx context.Context, cons jetstream.Consumer, subject string, handler messageHandler) error {
	msgs, err := cons.Fetch(1, jetstream.FetchMaxWait(c.Config.MaxWait))
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			g.Log().Errorf(ctx, "Error fetching messages for subject %s: %v", subject, err)
		}

		return err
	}

	return c.processFetchedMessages(ctx, msgs, handler, subject)
}

func (c *Client) processFetchedMessages(ctx context.Context, msgs jetstream.MessageBatch, handler messageHandler, subject string) error {
	for msg := range msgs.Messages() {
		if err := c.handleMessage(ctx, msg, handler); err != nil {
			g.Log().Errorf(ctx, "Error processing message: %v", err)
		}
	}

	if err := msgs.Error(); err != nil {
		g.Log().Errorf(ctx, "Error in message batch for subject %s: %v", subject, err)
		return err
	}

	return nil
}

func (c *Client) handleMessage(ctx context.Context, msg jetstream.Msg, handler messageHandler) error {
	err := handler(ctx, msg)
	if err == nil {
		if ackErr := msg.Ack(); ackErr != nil {
			g.Log().Errorf(ctx, "Error sending ACK for message: %v", ackErr)
			return ackErr
		}

		return nil
	}

	g.Log().Errorf(ctx, "Error handling message: %v", err)

	if nakErr := msg.Nak(); nakErr != nil {
		g.Log().Debugf(ctx, "Error sending NAK for message: %v", nakErr)

		return nakErr
	}

	return err
}

// Close closes the Client.
func (c *Client) Close(ctx context.Context) error {
	c.subManager.Close()

	if c.connManager != nil {
		c.connManager.Close(ctx)
	}

	return nil
}

// CreateTopic creates a new topic (stream) in NATS jStream.
func (c *Client) CreateTopic(ctx context.Context, name string) error {
	return c.streamManager.CreateStream(ctx, StreamConfig{
		Stream:   name,
		Subjects: []string{name},
	})
}

// DeleteTopic deletes a topic (stream) in NATS jStream.
func (c *Client) DeleteTopic(ctx context.Context, name string) error {
	return c.streamManager.DeleteStream(ctx, name)
}

// CreateStream creates a new stream in NATS jStream.
func (c *Client) CreateStream(ctx context.Context, cfg StreamConfig) error {
	return c.streamManager.CreateStream(ctx, cfg)
}

// DeleteStream deletes a stream in NATS jStream.
func (c *Client) DeleteStream(ctx context.Context, name string) error {
	return c.streamManager.DeleteStream(ctx, name)
}

// CreateOrUpdateStream creates or updates a stream in NATS jStream.
func (c *Client) CreateOrUpdateStream(ctx context.Context, cfg *jetstream.StreamConfig) (jetstream.Stream, error) {
	return c.streamManager.CreateOrUpdateStream(ctx, cfg)
}

// GetJetStreamStatus returns the status of the jStream connection.
func getJetStreamStatus(ctx context.Context, js jetstream.JetStream) (string, error) {
	_, err := js.AccountInfo(ctx)
	if err != nil {
		return jetStreamStatusError, err
	}

	return jetStreamStatusOK, nil
}
