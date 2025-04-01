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

const defaultRetryTimeout = 5 * time.Second

// Client represents a Client for NATS jStream operations.
type Client struct {
	natsOpts         []nats.Option
	connMgr          ConnMgr          // 连接到服务的连接管理器
	jsSubMgr         JsSubMgr         // jetstream的发布和消费
	subMgr           SubMgr           // 普通订阅发布
	serverAddr       string           // nats服务地址
	natsConnector    Connector        // 方便单元测试封装mock接口
	jetStreamCreator JetStreamCreator // 方便单元测试封装mock接口
}

// New 创建一个新的客户端
func New(srvAddr string, natsOpts ...nats.Option) *Client {
	c := &Client{
		serverAddr: srvAddr,
	}
	if len(natsOpts) == 0 {
		c.natsOpts = []nats.Option{
			nats.Name("Sieyuan NATS Client"),
		}
	}
	c.natsOpts = append(c.natsOpts, nats.NoEcho())
	return c
}

type clientOpts func(c *Client)

func WithJsManager(jsm JsSubMgr) func(c *Client) {
	return func(c *Client) {
		c.jsSubMgr = jsm
	}
}
func WithSubManager(subm SubMgr) func(c *Client) {
	return func(c *Client) {
		c.subMgr = subm
	}
}

// Connect establishes a connection to NATS and sets up jStream.
func (c *Client) Connect(ctx context.Context, opts ...clientOpts) error {
	if c.connMgr != nil && c.connMgr.Health().Status == health.StatusUp {
		g.Log().Warning(ctx, "NATS connection already established")
		return nil
	}
	for _, opt := range opts {
		opt(c)
	}

	c.connMgr = newConnMgr(c.serverAddr, c.natsConnector, c.jetStreamCreator)
	c.connMgr.Connect(ctx)

	return nil
}
func (c *Client) Conn() (*nats.Conn, error) {
	js, err := c.JetStream()
	if err != nil {
		return nil, err
	}
	return js.Conn(), nil
}
func (c *Client) JetStream() (jetstream.JetStream, error) {
	if c.connMgr == nil || c.connMgr.Health().Status == health.StatusDown {
		return nil, gerror.New("NATS connection not established")
	}
	js, err := c.connMgr.GetJetStream()
	if err != nil {
		return nil, err
	}
	return js, nil
}
func (c *Client) Flush() error {
	js, err := c.JetStream()
	if err != nil {
		return err
	}
	return js.Conn().Flush()
}

// Publish publishes a message to a topic.
func (c *Client) Publish(ctx context.Context, subject string, message []byte) error {
	return c.connMgr.Publish(ctx, subject, message)
}
func (c *Client) JsPublish(ctx context.Context, subject string, message []byte) error {
	return c.connMgr.JsPublish(ctx, subject, message)
}

// 流订阅
func (c *Client) JsSubscribe(ctx context.Context, stream streamIntf, identity []string, consumeType SubType, handler pubsub.SubscribeFunc) error {
	if !c.connMgr.isConnected() {
		return errClientNotConnected
	}
	err := c.jsSubMgr.NewSubscriber(ctx, stream, identity, consumeType)
	if err != nil {
		return err
	}
	if err := c.jsSubMgr.Subscribe(ctx, identity, handler); err != nil {
		return err
	}
	return nil
}

// 消息订阅
func (c *Client) Subscribe(ctx context.Context, topicName string, handler pubsub.SubscribeFunc) error {
	if !c.connMgr.isConnected() {
		return errClientNotConnected
	}
	conn, err := c.Conn()
	if err != nil {
		return err
	}
	if err := c.subMgr.NewSubscriber(ctx, conn, topicName, SubTypeSubAsync); err != nil {
		return err
	}
	if err := c.subMgr.Subscribe(ctx, topicName, handler); err != nil {
		return err
	}
	return nil
}

// Close closes the Client.
func (c *Client) Close(ctx context.Context) error {
	if c.subMgr != nil {
		c.subMgr.Close(ctx)
		c.subMgr = nil
	}
	if c.jsSubMgr != nil {
		c.jsSubMgr.Close(ctx)
		c.jsSubMgr = nil
	}
	if c.connMgr != nil {
		c.connMgr.Close(ctx)
		c.connMgr = nil
	}
	g.Log().Infof(ctx, "nats client close ok")
	return nil
}
