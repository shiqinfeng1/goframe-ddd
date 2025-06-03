package natsclient

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
)

//go:generate mockgen -destination=mock_jetstream.go -package=natsclient github.com/nats-io/nats.go/jetstream  KeyValue,ObjectStore,JetStream,KeyValueEntry

type Configs struct {
	Server string
	Bucket string
}

type Client struct {
	conn    *nats.Conn
	js      jetstream.JetStream
	kv      jetstream.KeyValue
	obj     jetstream.ObjectStore
	configs *Configs
}

// New creates a new NATS-KV client with the provided configuration.
func New(configs Configs) *Client {
	return &Client{configs: &configs}
}

func (c *Client) sendOperationStats(ctx context.Context, start time.Time, methodType string, kv ...string) {
	duration := time.Since(start)

	var key string
	if len(kv) > 0 {
		key = kv[0]
	}

	g.Log().Debug(ctx, &Log{
		Type:     methodType,
		Duration: duration.Microseconds(),
		Key:      key,
		Value:    c.configs.Bucket,
	})

	metrics.RecordHistogram(ctx, float64(duration.Milliseconds()),
		"bucket", c.configs.Bucket,
		"operation", methodType)
}
