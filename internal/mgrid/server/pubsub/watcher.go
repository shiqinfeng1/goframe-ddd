package pubsub

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

func (c *ControllerV1) startWatch(ctx context.Context, nc *nats.Conn) error {
	kvbkts := g.Cfg().MustGet(ctx, "nats.kvbuckets").Strings()
	objbkts := g.Cfg().MustGet(ctx, "nats.objbuckets").Strings()

	c.natsClient.StartWatch(ctx, nc, kvbkts, objbkts, nil)
	return nil
}
