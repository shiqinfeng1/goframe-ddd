package pubsub

import (
	"context"

	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

func (c *ControllerV1) startWatch(ctx context.Context, nc *natsclient.Conn) error {
	c.natsClient.StartWatch(ctx, nc, c.cfg.Nats.KvBuckets, c.cfg.Nats.ObjBuckets, nil)
	return nil
}
