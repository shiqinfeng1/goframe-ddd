package nats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
)

//go:generate mockgen -destination=mock_jetstream.go -package=nats github.com/nats-io/nats.go/jetstream  KeyValue,JetStream,KeyValueEntry

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

// Connect establishes a connection to NATS-KV and registers metrics using the provided configuration when the client is created.
func (c *Client) Connect(ctx context.Context) {
	g.Log().Debugf(ctx, "connecting to NATS-KV Store at %v with bucket %q", c.configs.Server, c.configs.Bucket)

	nc, err := nats.Connect(c.configs.Server)
	if err != nil {
		g.Log().Errorf(ctx, "error while connecting to NATS: %v", err)
		return
	}
	c.conn = nc
	g.Log().Debug(ctx, "connection to NATS successful")

	js, err := jetstream.New(nc)
	if err != nil {
		g.Log().Errorf(ctx, "error while initializing JetStream: %v", err)
		return
	}
	c.js = js
	g.Log().Debug(ctx, "jetStream initialized successfully")

	kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: c.configs.Bucket,
	})
	if err != nil {
		g.Log().Errorf(ctx, "error while creating/accessing KV bucket: %v", err)
		return
	}
	c.kv = kv
	g.Log().Infof(ctx, "successfully connected to NATS-KV Store at %s:%s ", c.configs.Server, c.configs.Bucket)

	obj, err := js.CreateObjectStore(ctx, jetstream.ObjectStoreConfig{
		Bucket: c.configs.Bucket,
	})
	if err != nil {
		g.Log().Errorf(ctx, "error while creating/accessing object bucket: %v", err)
		return
	}
	c.obj = obj
	g.Log().Infof(ctx, "successfully connected to NATS-object Store at %s:%s ", c.configs.Server, c.configs.Bucket)
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	defer c.sendOperationStats(ctx, time.Now(), "GET", "get", key)

	entry, err := c.kv.Get(ctx, key)
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			return "", fmt.Errorf("%w: %s", nats.ErrKeyNotFound, key)
		}

		return "", fmt.Errorf("failed to get key: %w", err)
	}

	return string(entry.Value()), nil
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET", "set", key, value)

	_, err := c.kv.Put(ctx, key, []byte(value))
	if err != nil {
		return fmt.Errorf("failed to set key-value pair: %w", err)
	}
	return nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	defer c.sendOperationStats(ctx, time.Now(), "DELETE", "delete", key)

	err := c.kv.Delete(ctx, key)
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			return fmt.Errorf("%w: %s", nats.ErrKeyNotFound, key)
		}
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}

func (c *Client) sendOperationStats(ctx context.Context, start time.Time, methodType, method string, kv ...string) {
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
