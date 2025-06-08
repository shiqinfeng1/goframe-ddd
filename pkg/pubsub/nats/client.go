package natsclient

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type Config struct {
	ServerUrl    string   `json:"serverUrl" yaml:"serverUrl"`
	ConsumerName string   `json:"consumerName" yaml:"consumerName"`
	StreamName   string   `json:"streamName" yaml:"streamName"`
	Subject1     string   `json:"subject1" yaml:"subject1"`
	Subject2     string   `json:"subject2" yaml:"subject2"`
	JSSubject1   string   `json:"jsSubject1" yaml:"jsSubject1"`
	JSSubject2   string   `json:"jsSubject2" yaml:"jsSubject2"`
	KVBuckets    []string `json:"kvBuckets" yaml:"kvBuckets"`
	ObjBuckets   []string `json:"objBuckets" yaml:"objBuckets"`
}

var defaultNatsOpts []nats.Option = []nats.Option{
	nats.NoEcho(),
}

// Client represents a Client for NATS jStream operations.
type Client struct {
	cfg        *Config
	logger     pubsub.Logger
	natsOpts   []nats.Option
	subscriber // 管理消息订阅者
	consumer   // 管理流的消费者
	watcher    // kv和object变化监听
}

// New 创建一个新的客户端
func New(cfg *Config, logger pubsub.Logger) *Client {
	c := &Client{
		cfg:    cfg,
		logger: logger,
		subscriber: subscriber{
			logger:        logger,
			subscriptions: gmap.NewStrAnyMap(true),
		},
		consumer: consumer{
			logger:        logger,
			subscriptions: gmap.NewStrAnyMap(true),
			exitNotify:    make(chan SubsKey),
		},
		watcher: watcher{
			logger:      logger,
			kvWatchers:  gmap.NewStrAnyMap(true),
			objWatchers: gmap.NewStrAnyMap(true),
		},
	}
	// 当订阅失败，或stream被删除后，需要删除相关资源
	go func() {
		ctx := gctx.New()
		for key := range c.consumer.exitNotify {
			c.consumer.Delete(ctx, key)
		}
		c.logger.Infof(ctx, "exit stream consumer ok")
	}()
	c.natsOpts = append(c.natsOpts, defaultNatsOpts...)
	return c
}

// 订阅消息
func (c *Client) SubMsg(ctx context.Context, nc *Conn, subject string, stype SubType, handler func(ctx context.Context, msg *nats.Msg) error) error {
	if err := c.subscriber.AddSubscription(ctx, nc, subject, stype, handler); err != nil {
		return err
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
func (c *Client) ConsumeStream(ctx context.Context, nc *Conn, sn, cn, subject string, stype SubType, handler func(ctx context.Context, msg *jetstream.Msg) error) error {
	js, err := c.jstream(nc)
	if err != nil {
		return err
	}
	cons, err := js.CreateOrUpdateConsumer(ctx, sn, cn, subject)
	if err != nil {
		return err
	}

	skey := NewSubsKey(subject, sn, cn)
	if err := c.consumer.Add(ctx, stype, skey, cons, handler, c.exitNotify); err != nil {
		return err
	}
	return nil
}

// 删除流消费
func (c *Client) DelConsumer(ctx context.Context, nc *Conn, sn, cn, subject string, stype SubType, handler ConsumeFunc) error {
	js, err := c.jstream(nc)
	if err != nil {
		return err
	}
	if err := js.DeleteConsumer(ctx, sn, cn, subject); err != nil {
		return err
	}
	skey := NewSubsKey(subject, sn, cn)
	if err := c.consumer.Delete(ctx, skey); err != nil {
		return err
	}
	return nil
}

// Close closes the Client.
func (c *Client) Close(ctx context.Context) error {
	if c == nil {
		return nil
	}
	if err := c.subscriber.Close(ctx); err != nil {
		return err
	}
	c.logger.Infof(ctx, "nats client close sub ok")
	if err := c.consumer.Close(ctx); err != nil {
		return err
	}
	c.logger.Infof(ctx, "nats client  close stream consume ok")
	if err := c.watcher.Stop(ctx); err != nil {
		return err
	}
	c.logger.Infof(ctx, "nats client close watcher ok")
	return nil
}
