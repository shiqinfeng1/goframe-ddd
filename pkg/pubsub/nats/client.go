package nats

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
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
	subManager       SubscriptionManagerIntf
	subscriptions    map[string]context.CancelFunc
	subMutex         sync.Mutex
	streamManager    StreamManagerIntf
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
		Config:        cfg,
		subManager:    newSubscriptionManager(),
		subscriptions: make(map[string]context.CancelFunc),
	}

	return client
}
func (c *Client) Conn() (*nats.Conn, error) {
	if c.connManager == nil || c.connManager.Health().Status == health.StatusDown {
		return nil, gerror.New("NATS connection not established")
	}
	js, err := c.connManager.getJetStream()
	if err != nil {
		return nil, err
	}
	return js.Conn(), nil
}

// Connect establishes a connection to NATS and sets up jStream.
func (c *Client) Connect(ctx context.Context) error {
	if c.connManager != nil && c.connManager.Health().Status == health.StatusUp {
		g.Log().Warning(ctx, "NATS connection already established")
		return nil
	}

	if err := c.validateAndPrepare(ctx); err != nil {
		return err
	}

	connManager := newConnectionManager(c.Config, c.natsConnector, c.jetStreamCreator)
	if err := connManager.Connect(ctx); err != nil {
		g.Log().Errorf(ctx, "failed to connect to NATS server at %v: %v", c.Config.Server, err)
		return err
	}

	c.connManager = connManager

	js, err := c.connManager.getJetStream()
	if err != nil {
		return err
	}

	c.streamManager = newStreamManager(js)
	c.subManager = newSubscriptionManager()

	return nil
}

func (c *Client) validateAndPrepare(ctx context.Context) error {
	if err := validateConfigs(c.Config); err != nil {
		g.Log().Errorf(ctx, "could not initialize NATS jStream: %v", err)

		return err
	}

	return nil
}

// Publish publishes a message to a topic.
func (c *Client) Publish(ctx context.Context, subject string, message []byte) error {
	return c.connManager.Publish(ctx, subject, message)
}

// Subscribe subscribes to a topic and returns a single message.
func (c *Client) Subscribe(ctx context.Context, topic string, handler pubsub.SubscribeFunc) error {
	if !c.connManager.isConnected() {
		time.Sleep(defaultRetryTimeout)
		return errClientNotConnected
	}

	js, err := c.connManager.getJetStream()
	if err != nil {
		return err
	}
	if err := c.subManager.Subscribe(ctx, topic, js, c.Config, handler); err != nil {
		return err
	}
	return nil
}

// 根据消费主题自动生成一个消费者的名字，带有通配符的主题，需要替换通配符
func generateConsumerName(consumer, subject string) string {
	subject = strings.ReplaceAll(subject, ".", "_")
	subject = strings.ReplaceAll(subject, "*", "token")
	subject = strings.ReplaceAll(subject, ">", "tokens")
	return fmt.Sprintf("%s_%s", consumer, subject)
}

func (c *Client) SubscribeWithHandler(ctx context.Context, subject string, handler messageHandler) error {
	c.subMutex.Lock()
	defer c.subMutex.Unlock()

	// Cancel any existing subscription for this subject
	c.cancelExistingSubscription(subject)

	js, err := c.connManager.getJetStream()
	if err != nil {
		return err
	}

	consumerName := generateConsumerName(c.Config.ConsumerName, subject)

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

// client 直接创建一个stream的comsumer
func (c *Client) createOrUpdateConsumer(
	ctx context.Context, js jetstream.JetStream, subject, consumerName string,
) (jetstream.Consumer, error) {
	cons, err := js.CreateOrUpdateConsumer(ctx, c.Config.Stream.Name, jetstream.ConsumerConfig{
		Durable:       consumerName,
		AckPolicy:     jetstream.AckNonePolicy, //AckExplicitPolicy,
		FilterSubject: subject,
		MaxDeliver:    c.Config.Stream.MaxDeliver,
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		g.Log().Errorf(ctx, "failed to create or update consumer for stream %v: %v", c.Config.Stream.Name, err)
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

// 直接获取消息
func (c *Client) fetchAndProcessMessages(ctx context.Context, cons jetstream.Consumer, subject string, handler messageHandler) error {
	msgs, err := cons.Messages(jetstream.PullMaxMessages(1))
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			g.Log().Errorf(ctx, "Error fetching messages for subject %s: %v", subject, err)
		}

		return err
	}

	return c.processFetchedMessages(ctx, msgs, handler, subject)
}

func (c *Client) processFetchedMessages(ctx context.Context, msgs jetstream.MessagesContext, handler messageHandler, subject string) error {
	g.Log().Infof(ctx, "ready to consume msg for %v", subject)
	for {
		msg, err := msgs.Next()
		if err != nil {
			g.Log().Warningf(ctx, "Error processing message subject %v: %v", subject, err)
			return nil
		}
		if err := c.handleMessage(ctx, msg, handler); err != nil {
			return err
		}
	}
}

func (c *Client) handleMessage(ctx context.Context, msg jetstream.Msg, handler messageHandler) error {
	err := handler(ctx, msg)
	if err != nil {
		g.Log().Errorf(ctx, "Error handling message: %v", err)
		if nakErr := msg.Nak(); nakErr != nil {
			g.Log().Debugf(ctx, "Error sending NAK for message: %v", nakErr)
			return nakErr
		}
		return err
	}

	if ackErr := msg.Ack(); ackErr != nil {
		g.Log().Errorf(ctx, "Error sending ACK for message: %v", ackErr)
		return ackErr
	}
	return nil
}

// Close closes the Client.
func (c *Client) Close(ctx context.Context) error {
	if c.subManager != nil {
		c.subManager.Close()
		c.subManager = nil
	}
	if c.connManager != nil {
		c.connManager.Close(ctx)
		c.connManager = nil
	}
	g.Log().Infof(ctx, "nats client close ok")
	return nil
}

// CreateTopic creates a new topic (stream) in NATS jStream.
func (c *Client) CreateTopic(ctx context.Context) error {
	return c.streamManager.CreateStream(ctx, c.Config.Stream)
}

// DeleteTopic deletes a topic (stream) in NATS jStream.
func (c *Client) DeleteTopic(ctx context.Context, name string) error {
	return c.streamManager.DeleteStream(ctx, name)
}

// GetJetStreamStatus returns the status of the jStream connection.
func getJetStreamStatus(ctx context.Context, js jetstream.JetStream) (string, error) {
	_, err := js.AccountInfo(ctx)
	if err != nil {
		return jetStreamStatusError, err
	}

	return jetStreamStatusOK, nil
}
