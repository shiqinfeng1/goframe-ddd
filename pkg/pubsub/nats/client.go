package nats

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

var defaultNatsOpts []nats.Option = []nats.Option{
	nats.NoEcho(),
}

// Client represents a Client for NATS jStream operations.
type Client struct {
	logger         pubsub.Logger
	natsOpts       []nats.Option
	connFactory    ConnFactory     // 连接到服务的连接管理器
	serverAddr     string          // nats服务地址
	subscriber     *Subscriber     // 管理消息订阅者
	streamConsumer *StreamConsumer // 管理流的消费者
}

// New 创建一个新的客户端
func New(ctx context.Context, logger pubsub.Logger, srvAddr string, natsConnector Connector) *Client {
	c := &Client{
		logger:         logger,
		serverAddr:     srvAddr,
		connFactory:    NewFactory(logger, srvAddr, natsConnector),
		subscriber:     NewSubscriber(logger),
		streamConsumer: NewStreamConsumer(logger),
	}
	c.natsOpts = append(c.natsOpts, defaultNatsOpts...)
	return c
}

// 获取一个新的连接，到nats服务端
func (c *Client) NewConn(ctx context.Context, opts ...nats.Option) (*Conn, error) {
	return c.connFactory.New(ctx, append(opts, defaultNatsOpts...)...)
}

// 订阅消息
func (c *Client) SubMsg(ctx context.Context, natsConn *Conn, subject string, stype SubType, handler SubscribeFunc) error {
	if err := c.subscriber.AddSubscription(ctx, natsConn, subject, stype, handler); err != nil {
		return gerror.Wrap(err, "add subscription fail")
	}
	return nil
}

// 创建流
func (c *Client) CreateStream(ctx context.Context, natsConn *Conn, streamName string, subjects []string) error {
	js, err := natsConn.JetStream()
	if err != nil {
		return err
	}
	jswrapper := NewJetStreamWrapper(c.logger, js)
	if err := jswrapper.CreateStream(ctx, streamName, subjects); err != nil {
		return err
	}
	return nil
}

// 创建或更新流
func (c *Client) CreateOrUpdateStream(ctx context.Context, natsConn *Conn, streamName string, subjects []string) error {
	js, err := natsConn.JetStream()
	if err != nil {
		return err
	}
	jswrapper := NewJetStreamWrapper(c.logger, js)
	if err := jswrapper.CreateOrUpdateStream(ctx, streamName, subjects); err != nil {
		return err
	}
	return nil
}

// 删除流,  该流上的所有consumer也会被自动删除
func (c *Client) DeleteStream(ctx context.Context, natsConn *Conn, streamName string) error {
	js, err := natsConn.JetStream()
	if err != nil {
		return err
	}
	jswrapper := NewJetStreamWrapper(c.logger, js)
	if err := jswrapper.DeleteStream(ctx, streamName); err != nil {
		return err
	}
	return nil
}

// 流消费
func (c *Client) ConsumeStream(ctx context.Context, natsConn *Conn, streamName, consumerName, subject string, stype SubType, handler ConsumeFunc) error {
	js, err := natsConn.JetStream()
	if err != nil {
		return err
	}
	// 在服务端创建消费者
	jswrapper := NewJetStreamWrapper(c.logger, js)
	consumer, err := jswrapper.CreateConsumer(ctx, streamName, consumerName, subject)
	if err != nil {
		return err
	}
	//
	skey := NewSubsKey(subject, streamName, consumerName)
	if err := c.streamConsumer.AddConsume(ctx, stype, skey, consumer, handler, c.streamConsumer.exitNotify); err != nil {
		return gerror.Wrap(err, "add consume fail")
	}
	return nil
}

// 删除流消费
func (c *Client) DeleteConsumer(ctx context.Context, natsConn *Conn, streamName, consumerName, subject string, stype SubType, handler ConsumeFunc) error {
	js, err := natsConn.JetStream()
	if err != nil {
		return err
	}
	jswrapper := NewJetStreamWrapper(c.logger, js)
	if err := jswrapper.DeleteConsumer(ctx, streamName, consumerName, subject); err != nil {
		return err
	}
	skey := NewSubsKey(subject, streamName, consumerName)
	if err := c.streamConsumer.DeleteConsumer(ctx, skey); err != nil {
		return gerror.Wrap(err, "add consume fail")
	}
	return nil
}

// Close closes the Client.
func (c *Client) Close(ctx context.Context) error {
	if c.subscriber != nil {
		c.subscriber.Close(ctx)
		c.subscriber = nil
		g.Log().Infof(ctx, "nats client '%v' close sub ok", c.serverAddr)
	}
	if c.streamConsumer != nil {
		c.streamConsumer.Close(ctx)
		c.streamConsumer = nil
		g.Log().Infof(ctx, "nats client '%v' close stream consume ok", c.serverAddr)
	}

	g.Log().Infof(ctx, "nats client '%v' close all ok", c.serverAddr)
	return nil
}
