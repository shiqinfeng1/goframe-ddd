package nats

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
)

//go:generate mockgen -destination=mock_jetstream.go -package=nats github.com/nats-io/nats.go/jetstream JetStream,Stream,Consumer,Msg,MessageBatch

type Conn struct {
	conn             ConnIntf
	jStream          jetstream.JetStream
	serverAddr       string // 服务地址
	natsConnector    Connector
	jetStreamCreator JetStreamCreator
}

func (c *Conn) GetJetStream() (jetstream.JetStream, error) {
	if c.jStream == nil {
		return nil, errJetStreamNotConfigured
	}

	return c.jStream, nil
}

// newConn creates a new Conn.
func newConn(
	serverAddr string,
	natsConnector Connector,
	jetStreamCreator JetStreamCreator,
) *Conn {
	// 设置连接器
	if natsConnector == nil {
		natsConnector = &defaultConnector{}
	}
	// 设置js构造器
	if jetStreamCreator == nil {
		jetStreamCreator = &defaultJetStreamCreator{}
	}

	return &Conn{
		natsConnector:    natsConnector,
		jetStreamCreator: jetStreamCreator,
		serverAddr:       serverAddr,
	}
}

// Connect establishes a connection to NATS and sets up JetStream.
// 异步重试连接
func (c *Conn) Connect(ctx context.Context, opts ...nats.Option) {
	for {
		conn, err := c.natsConnector.Connect(c.serverAddr, opts...)
		if err != nil {
			g.Log().Warningf(ctx, "try to connect to NATS server at %v: %v", c.serverAddr, err)
			time.Sleep(defaultRetryTimeout)
			continue
		}
		// 连接成功后，创建jetstream
		js, err := c.jetStreamCreator.New(conn)
		if err != nil {
			conn.Close()
			g.Log().Debugf(ctx, "Failed to create jStream context: %v", err)
			time.Sleep(defaultRetryTimeout)
			continue
		}

		c.conn = conn
		c.jStream = js
		g.Log().Infof(ctx, "Successfully connected to NATS server at %v by '%v'", c.serverAddr, conn.Conn().Opts.Name)
		return
	}
}

func (c *Conn) Close(_ context.Context) {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Conn) Subscribe(ctx context.Context, subject string, handler func(msg *nats.Msg)) error {
	metrics.IncrementCounter(ctx, metrics.NatsPublishTotalCount, "subject", subject)
	subs, err := c.conn.Conn().Subscribe(subject, handler)
	if err != nil {
		return err
	}
	subs.Unsubscribe()
	metrics.IncrementCounter(ctx, metrics.NatsPublishSuccessCount, "subject", subject)
	return nil
}
func (c *Conn) Publish(ctx context.Context, subject string, message []byte) error {
	metrics.IncrementCounter(ctx, metrics.NatsPublishTotalCount, "subject", subject)
	if err := c.conn.Conn().Publish(subject, message); err != nil {
		return err
	}
	metrics.IncrementCounter(ctx, metrics.NatsPublishSuccessCount, "subject", subject)
	return nil
}
func (c *Conn) JsPublish(ctx context.Context, subject string, message []byte) error {
	metrics.IncrementCounter(ctx, metrics.NatsJsPublishTotalCount, "subject", subject)

	if err := c.validateJetStream(ctx, subject); err != nil {
		return err
	}
	// 发布消息
	_, err := c.jStream.Publish(ctx, subject, message)
	if err != nil {
		g.Log().Errorf(ctx, "failed to publish message to NATS jStream: %v", err)
		return err
	}

	metrics.IncrementCounter(ctx, metrics.NatsJsPublishSuccessCount, "subject", subject)

	return nil
}

func (c *Conn) validateJetStream(_ context.Context, subject string) error {
	if c.jStream == nil || subject == "" {
		err := errJetStreamNotConfigured
		return err
	}

	return nil
}

// 返回nats客户端和服务端之间连接的健康状态
func (c *Conn) Health() *health.Health {
	if c.conn == nil {
		return &health.Health{
			Status: health.StatusDown,
		}
	}

	status := c.conn.Status()
	if status == nats.CONNECTED {
		return &health.Health{
			Status: health.StatusUp,
			Details: map[string]any{
				"server": c.serverAddr,
			},
		}
	}

	return &health.Health{
		Status: health.StatusDown,
		Details: map[string]any{
			"server": c.serverAddr,
		},
	}
}

func (c *Conn) isConnected() bool {
	if c.conn == nil {
		return false
	}

	return c.conn.Status() == nats.CONNECTED
}
