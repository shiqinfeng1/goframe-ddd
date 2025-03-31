package nats

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// natsCommitter implements the pubsub.Committer interface for Client messages.
type natsCommitter struct {
	msg *nats.Msg
}

// Commit commits the message.
func (c *natsCommitter) Commit(ctx context.Context) {
	if err := c.msg.Ack(); err != nil {
		g.Log().Errorf(ctx, "Error committing message:%v", err)

		// nak the message
		if err := c.msg.Nak(); err != nil {
			g.Log().Errorf(ctx, "Error naking message:%v", err)
		}
		return
	}
}

type jsCommitter struct {
	msg jetstream.Msg
}

// Commit commits the message.
func (c *jsCommitter) Commit(ctx context.Context) {
	if err := c.msg.Ack(); err != nil {
		g.Log().Errorf(ctx, "Error committing message:%v", err)

		// nak the message
		if err := c.msg.Nak(); err != nil {
			g.Log().Errorf(ctx, "Error naking message:%v", err)
		}
		return
	}
}
