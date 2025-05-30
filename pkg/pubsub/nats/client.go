package nats

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

var defaultNatsOpts []nats.Option = []nats.Option{
	nats.NoEcho(),
}

// Client represents a Client for NATS jStream operations.
type Client struct {
	logger     pubsub.Logger
	natsOpts   []nats.Option
	subscriber // 管理消息订阅者
	consumer   // 管理流的消费者
	watcher    // kv和object变化监听
}

// New 创建一个新的客户端
func New(logger pubsub.Logger) *Client {
	c := &Client{
		logger: logger,
		subscriber: subscriber{
			logger:        logger,
			subscriptions: make(map[string]*subscription),
			subMutex:      sync.Mutex{},
		},
		consumer: consumer{
			logger:        logger,
			subscriptions: make(map[string]*streamConsume),
			subMutex:      sync.Mutex{},
			exitNotify:    make(chan SubsKey),
		},
		watcher: watcher{
			logger:      logger,
			kvWatchers:  make(map[string]jetstream.KeyWatcher),
			objWatchers: make(map[string]jetstream.ObjectWatcher),
		},
	}
	// 当订阅失败，或stream被删除后，需要删除相关资源
	go func() {
		ctx := gctx.New()
		for key := range c.consumer.exitNotify {
			c.consumer.Delete(ctx, key)
		}
		c.logger.Debugf(ctx, "exit stream consumer ok")
	}()
	c.natsOpts = append(c.natsOpts, defaultNatsOpts...)
	return c
}

// 订阅消息
func (c *Client) SubMsg(ctx context.Context, nc *Conn, subject string, stype SubType, handler SubscribeFunc) error {
	if err := c.subscriber.AddSubscription(ctx, nc, subject, stype, handler); err != nil {
		return gerror.Wrap(err, "add subscription fail")
	}
	return nil
}
func (c *Client) jstream(nc *Conn) (*JetStreamWrapper, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	return NewJetStreamWrapper(c.logger, js), nil
}

// 创建流
func (c *Client) CreateStream(ctx context.Context, nc *Conn, streamName string, subjects []string) error {
	js, err := c.jstream(nc)
	if err != nil {
		return err
	}
	if err := js.CreateStream(ctx, streamName, subjects); err != nil {
		return err
	}
	return nil
}

// 创建或更新流
func (c *Client) CreateOrUpdateStream(ctx context.Context, nc *Conn, streamName string, subjects []string) error {
	js, err := c.jstream(nc)
	if err != nil {
		return err
	}
	if err := js.CreateOrUpdateStream(ctx, streamName, subjects); err != nil {
		return err
	}
	return nil
}

// 删除流,  该流上的所有consumer也会被自动删除
func (c *Client) DeleteStream(ctx context.Context, nc *Conn, streamName string) error {
	js, err := c.jstream(nc)
	if err != nil {
		return err
	}
	if err := js.DeleteStream(ctx, streamName); err != nil {
		return err
	}
	return nil
}

// 流消费
func (c *Client) ConsumeStream(ctx context.Context, nc *Conn, streamName, consumerName, subject string, stype SubType, handler ConsumeFunc) error {
	js, err := c.jstream(nc)
	if err != nil {
		return err
	}
	cons, err := js.CreateOrUpdateConsumer(ctx, streamName, consumerName, subject)
	if err != nil {
		return err
	}
	//
	skey := NewSubsKey(subject, streamName, consumerName)
	if err := c.consumer.Add(ctx, stype, skey, cons, handler, c.exitNotify); err != nil {
		return gerror.Wrap(err, "add consume fail")
	}
	return nil
}

// 删除流消费
func (c *Client) DelConsumer(ctx context.Context, nc *Conn, streamName, consumerName, subject string, stype SubType, handler ConsumeFunc) error {
	js, err := c.jstream(nc)
	if err != nil {
		return err
	}
	if err := js.DeleteConsumer(ctx, streamName, consumerName, subject); err != nil {
		return err
	}
	skey := NewSubsKey(subject, streamName, consumerName)
	if err := c.consumer.Delete(ctx, skey); err != nil {
		return gerror.Wrap(err, "add consume fail")
	}
	return nil
}

// Close closes the Client.
func (c *Client) Close(ctx context.Context) error {
	if err := c.subscriber.Close(ctx); err != nil {
		return err
	}
	g.Log().Infof(ctx, "nats client close sub ok")
	if err := c.consumer.Close(ctx); err != nil {
		return err
	}
	g.Log().Infof(ctx, "nats client  close stream consume ok")
	if err := c.watcher.Stop(); err != nil {
		return err
	}
	g.Log().Infof(ctx, "nats client close watcher ok")
	return nil
}
