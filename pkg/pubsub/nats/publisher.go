package natsclient

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
)

type Publisher struct {
	conn      *nats.Conn
	jetstream jetstream.JetStream
}

func NewPublisher(conn *nats.Conn) *Publisher {
	jetstream, _ := jetstream.New(conn)
	return &Publisher{
		conn:      conn,
		jetstream: jetstream,
	}
}
func (c *Publisher) JetStream() jetstream.JetStream {
	return c.jetstream
}

func (c *Publisher) Close(_ context.Context) {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Publisher) PublishMsg(ctx context.Context, subject string, message []byte) error {
	metrics.Inc(ctx, metrics.NatsPublishTotalCount, "topic", subject)
	if err := c.conn.Publish(subject, message); err != nil {
		return gerror.Wrapf(err, "publish msg fail. subject=%v", subject)
	}
	metrics.Inc(ctx, metrics.NatsPublishSuccessCount, "topic", subject)
	return nil
}
func (c *Publisher) PublishStreamMsg(ctx context.Context, subject string, message []byte) error {
	metrics.Inc(ctx, metrics.NatsJSPublishTotalCount, "topic", subject)
	_, err := c.jetstream.Publish(ctx, subject, message)
	if err != nil {
		return gerror.Wrapf(err, "publish stream msg fail. subject=%v", subject)
	}
	metrics.Inc(ctx, metrics.NatsJSPublishSuccessCount, "topic", subject)
	return nil
}
