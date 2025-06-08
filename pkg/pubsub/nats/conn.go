package natsclient

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
)

const (
	defaultRetryTimeout = 3 * time.Second
)

type Conn struct {
	conn       natsConn
	jStream    jetstream.JetStream
	serverAddr string
}

func (c *Conn) JetStream() (jetstream.JetStream, error) {
	if !c.conn.IsConnected() {
		return nil, gerror.New("nats not connected")
	}
	if c.jStream != nil {
		return c.jStream, nil
	}
	c.jStream, _ = jetstream.New(c.conn.NatsConn())
	return c.jStream, nil
}

func (c *Conn) Close(_ context.Context) {
	if c != nil && c.conn != nil {
		c.conn.Close()
	}
}

func (c *Conn) SubMsg(ctx context.Context, subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	if err := c.validateConn(ctx, subject); err != nil {
		return nil, err
	}
	subs, err := c.conn.NatsConn().Subscribe(subject, handler)
	if err != nil {
		return nil, gerror.Wrapf(err, "subscribe subject=%v fail", subject)
	}
	return subs, nil
}
func (c *Conn) PubMsg(ctx context.Context, subject string, message []byte) error {
	metrics.Inc(ctx, metrics.NatsPublishTotalCount, "subject", subject)
	if err := c.validateConn(ctx, subject); err != nil {
		return err
	}
	if err := c.conn.NatsConn().Publish(subject, message); err != nil {
		return gerror.Wrapf(err, "publish subject=%v fail", subject)
	}
	metrics.Inc(ctx, metrics.NatsPublishSuccessCount, "subject", subject)
	return nil
}
func (c *Conn) PubStream(ctx context.Context, subject string, message []byte) error {
	metrics.Inc(ctx, metrics.NatsJsPublishTotalCount, "subject", subject)
	if err := c.validateJetStream(ctx, subject); err != nil {
		return err
	}
	// 发布消息
	_, err := c.jStream.Publish(ctx, subject, message)
	if err != nil {
		return err
	}
	metrics.Inc(ctx, metrics.NatsJsPublishSuccessCount, "subject", subject)

	return nil
}
func (c *Conn) validateConn(_ context.Context, subject string) error {
	if !c.conn.IsConnected() {
		return gerror.New("nats not connected")
	}
	if subject == "" {
		return gerror.New("subject is nil")
	}
	return nil
}
func (c *Conn) validateJetStream(ctx context.Context, subject string) error {
	if err := c.validateConn(ctx, subject); err != nil {
		return err
	}
	if c.jStream == nil {
		return gerror.New("jetstream not configed")
	}
	return nil
}
