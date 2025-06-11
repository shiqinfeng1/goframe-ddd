package natsclient

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type Config struct {
	ServerUrl     string   `json:"serverUrl" yaml:"serverUrl"`
	StreamName    string   `json:"streamName" yaml:"streamName"`
	Subject1      string   `json:"subject1" yaml:"subject1"`
	Subject2      string   `json:"subject2" yaml:"subject2"`
	JSSubject1    string   `json:"jsSubject1" yaml:"jsSubject1"`
	ConsumerName1 string   `json:"consumerName1" yaml:"consumerName1"`
	JSSubject2    string   `json:"jsSubject2" yaml:"jsSubject2"`
	ConsumerName2 string   `json:"consumerName2" yaml:"consumerName2"`
	KVBuckets     []string `json:"kvBuckets" yaml:"kvBuckets"`
	ObjBuckets    []string `json:"objBuckets" yaml:"objBuckets"`
}

// Client represents a Client for NATS jStream operations.
type Client struct {
	cfg            *Config
	logger         pubsub.Logger
	AsyncSubscribe *AsyncSubscriber
	SyncSubscribe  *SyncSubscriber
	SyncConsumer   *SyncConsumer
	watcher        // kv和object变化监听
}

// New 创建一个新的客户端
func New(cfg *Config, logger pubsub.Logger) *Client {
	c := &Client{
		cfg:    cfg,
		logger: logger,
		watcher: watcher{
			logger:      logger,
			kvWatchers:  gmap.NewStrAnyMap(true),
			objWatchers: gmap.NewStrAnyMap(true),
		},
	}
	return c
}

// 异步订阅消息
func (c *Client) NewAsyncSubscriber(logger pubsub.Logger, f Factory) (*AsyncSubscriber, error) {
	if c.AsyncSubscribe != nil {
		return c.AsyncSubscribe, nil
	}
	asub, err := NewAsyncSubscriber(logger, f)
	if err != nil {
		return nil, err
	}
	c.AsyncSubscribe = asub
	return c.AsyncSubscribe, nil
}
func (c *Client) NewSyncSubscriber(logger pubsub.Logger, f Factory) (*SyncSubscriber, error) {
	if c.SyncSubscribe != nil {
		return c.SyncSubscribe, nil
	}
	asub, err := NewSyncSubscriber(logger, f)
	if err != nil {
		return nil, err
	}
	c.SyncSubscribe = asub
	return c.SyncSubscribe, nil
}
func (c *Client) NewSyncConsumer(logger pubsub.Logger, f Factory) (*SyncConsumer, error) {
	if c.SyncConsumer != nil {
		return c.SyncConsumer, nil
	}
	asub, err := NewSyncConsumer(logger, f)
	if err != nil {
		return nil, err
	}
	c.SyncConsumer = asub
	return c.SyncConsumer, nil
}

// 取消异步订阅
func (c *Client) ShutdownSubscribe() {
	if c.AsyncSubscribe != nil {
		c.AsyncSubscribe.Shutdown()
	}
	if c.SyncSubscribe != nil {
		c.SyncSubscribe.Shutdown()
	}
	if c.SyncConsumer != nil {
		c.SyncConsumer.Shutdown()
	}
}

// 创建或更新流
func (c *Client) CreateOrUpdateStream(ctx context.Context, nc *nats.Conn, streamName string, subjects []string) error {
	jetstream, _ := jetstream.New(nc)
	js := NewJetStreamWrapper(c.logger, jetstream)
	if err := js.CreateOrUpdateStream(ctx, streamName, subjects); err != nil {
		return err
	}
	return nil
}

// 删除流,  该流上的所有consumer也会被自动删除
func (c *Client) DeleteStream(ctx context.Context, nc *nats.Conn, streamName string) error {
	jetstream, _ := jetstream.New(nc)
	js := NewJetStreamWrapper(c.logger, jetstream)
	if err := js.DeleteStream(ctx, streamName); err != nil {
		return err
	}
	return nil
}

// Close closes the Client.
func (c *Client) Close(ctx context.Context) error {
	if c == nil {
		return nil
	}
	c.ShutdownSubscribe()
	c.logger.Infof(ctx, "nats client close subscribe/consumer ok")

	if err := c.watcher.StopAllWatch(ctx); err != nil {
		return err
	}
	c.logger.Infof(ctx, "nats client close watcher ok")
	return nil
}
