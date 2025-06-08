package natsclient

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Connect establishes a connection to NATS-KV and registers metrics using the provided configuration when the client is created.
func (c *Client) Connect(ctx context.Context) {
	g.Log().Debugf(ctx, "connecting to NATS-KV/OBJ Store at %v", c.cfg.Server)

	nc, err := nats.Connect(c.cfg.Server)
	if err != nil {
		c.logger.Errorf(ctx, "error while connecting to NATS: %v", err)
		return
	}
	c.conn = nc
	c.logger.Debugf(ctx, "connection to NATS successful")

	js, err := jetstream.New(nc)
	if err != nil {
		c.logger.Errorf(ctx, "error while initializing JetStream: %v", err)
		return
	}
	c.js = js
	c.logger.Debugf(ctx, "jetStream initialized successfully")
	for _, bucket := range c.cfg.KVBuckets {
		kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
			Bucket: bucket,
		})
		if err != nil {
			c.logger.Errorf(ctx, "error while creating/accessing KV bucket: %v", err)
			return
		}
		if notexist := c.kv.SetIfNotExist(bucket, kv); !notexist {
			c.logger.Errorf(ctx, "KV bucket <%v> already exists", bucket)
			return
		}
	}
	c.logger.Infof(ctx, "successfully connected to NATS-KV Store at %s:%s ", c.cfg.Server, c.cfg.KVBuckets)
	for _, bucket := range c.cfg.ObjBuckets {
		obj, err := js.CreateOrUpdateObjectStore(ctx, jetstream.ObjectStoreConfig{
			Bucket: bucket,
		})
		if err != nil {
			c.logger.Errorf(ctx, "error while creating/accessing KV bucket: %v", err)
			return
		}
		if notexist := c.obj.SetIfNotExist(bucket, obj); !notexist {
			c.logger.Errorf(ctx, "OBJ bucket <%v> already exists", bucket)
			return
		}
	}
	c.logger.Infof(ctx, "successfully connected to NATS-object Store at %s:%s ", c.cfg.Server, c.cfg.ObjBuckets)
}

func (c *Client) Close(ctx context.Context) {
	c.conn.Close()
}
