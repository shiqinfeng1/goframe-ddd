package nats

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

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

	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: c.configs.Bucket,
	})
	if err != nil {
		g.Log().Errorf(ctx, "error while creating/accessing KV bucket: %v", err)
		return
	}
	c.kv = kv
	g.Log().Infof(ctx, "successfully connected to NATS-KV Store at %s:%s ", c.configs.Server, c.configs.Bucket)

	obj, err := js.CreateOrUpdateObjectStore(ctx, jetstream.ObjectStoreConfig{
		Bucket: c.configs.Bucket,
	})
	if err != nil {
		g.Log().Errorf(ctx, "error while creating/accessing object bucket: %v", err)
		return
	}
	c.obj = obj
	g.Log().Infof(ctx, "successfully connected to NATS-object Store at %s:%s ", c.configs.Server, c.configs.Bucket)
}

func (c *Client) Close(ctx context.Context) {
	c.conn.Close()
}
