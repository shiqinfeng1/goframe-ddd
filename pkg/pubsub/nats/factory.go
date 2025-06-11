package natsclient

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type Factory interface {
	New(ctx context.Context, name string, opts ...nats.Option) (*nats.Conn, error)
}
type factory struct {
	logger     pubsub.Logger
	serverAddr string
}

func NewFactory(
	logger pubsub.Logger,
) Factory {
	return &factory{
		logger:     logger,
		serverAddr: g.Cfg().MustGet(gctx.New(), "nats.serverUrl").String(),
	}
}
func (f *factory) New(ctx context.Context, name string, opts ...nats.Option) (*nats.Conn, error) {
	opts = append(opts,
		nats.Name(name),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			f.logger.Infof(ctx, "nats client disconnect CB: client '%v' disconnected: %v", name, err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			f.logger.Infof(ctx, "nats client reconnect CB: '%v' reconnected", name)
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			f.logger.Infof(ctx, "nats client close CB: '%v' closed", name)
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			f.logger.Infof(ctx, "nats client error CB: '%v' occur error: %v", name, err)
		}),
	)

	conn, err := nats.Connect(f.serverAddr, opts...)
	if err != nil {
		return nil, err
	}

	if !conn.IsConnected() {
		return nil, gerror.New("connect to nats timeout")
	}

	f.logger.Infof(ctx, "successfully connected to NATS server at %v by '%v'", f.serverAddr, conn.Opts.Name)
	return conn, nil
}
