package pubsub

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

func (c *ControllerV1) startWatch(ctx context.Context, nc *nats.Conn) error {

	js, _ := jetstream.New(nc)

	for _, bkt := range c.cfg.Nats.KVBuckets {
		_, err := js.KeyValue(ctx, bkt)
		if gerror.Is(err, jetstream.ErrBucketNotFound) {
			js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
				Bucket: bkt,
			})
		}
	}
	// 启动对象存储监听
	for _, bkt := range c.cfg.Nats.ObjBuckets {
		_, err := js.ObjectStore(ctx, bkt)
		if gerror.Is(err, jetstream.ErrBucketNotFound) {
			js.CreateObjectStore(ctx, jetstream.ObjectStoreConfig{
				Bucket: bkt,
			})
		}
	}
	publisher := natsclient.NewPublisher(nc)
	return c.natsClient.StartWatch(ctx, publisher, c.cfg.Nats.KVBuckets, c.cfg.Nats.ObjBuckets, nil)
}
