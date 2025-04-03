package nats

// // jsCommitter implements the pubsub.Committer interface for Client messages.
// type jsCommitter struct {
// 	msg jetstream.Msg
// }

// // Commit commits the message.
// func (c *jsCommitter) Commit(ctx context.Context) {
// 	if err := c.msg.Ack(); err != nil {
// 		g.Log().Errorf(ctx, "Error committing message:%v", err)

// 		// nak the message
// 		if err := c.msg.Nak(); err != nil {
// 			g.Log().Errorf(ctx, "Error naking message:%v", err)
// 		}
// 		return
// 	}
// }
