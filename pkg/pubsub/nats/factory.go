package natsclient

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type Factory interface {
	New(ctx context.Context, name string, opts ...nats.Option) (*Conn, error)
}
type factory struct {
	logger     pubsub.Logger
	serverAddr string
	connector  Connector
}

func NewFactory(
	logger pubsub.Logger,
	natsConnector Connector,
) Factory {
	// 设置连接器
	if natsConnector == nil {
		natsConnector = &defaultConnector{}
	}

	return &factory{
		logger:     logger,
		connector:  natsConnector,
		serverAddr: g.Cfg().MustGet(gctx.New(), "nats.serverUrl").String(),
	}
}
func (f *factory) New(ctx context.Context, name string, opts ...nats.Option) (*Conn, error) {
	opts = append(opts,
		nats.Name(name),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			f.logger.Infof(ctx, "nats client '%v' disconnected: %v", name, err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			f.logger.Infof(ctx, "nats client '%v' reconnected", name)
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			f.logger.Infof(ctx, "nats client '%v' closed", name)
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			f.logger.Infof(ctx, "nats client '%v' occur error: %v", name, err)
		}),
	)

	conn, err := f.connector.Connect(f.serverAddr, opts...)
	for i := 0; i < defaultRetryCount && !conn.IsConnected(); i++ {
		f.logger.Warningf(ctx, "[%v/%v]try to connect to NATS server at %v: %v", i+1, defaultRetryCount, f.serverAddr, err)
		time.Sleep(defaultRetryTimeout)
	}
	if !conn.IsConnected() {
		return nil, gerror.New("connect to nats timeout")
	}

	// 连接成功后，创建jetstream
	js, err := conn.NewJetStream()
	if err != nil {
		conn.Close()
		return nil, gerror.Wrap(err, "failed to create jStream context")
	}
	f.logger.Infof(ctx, "successfully connected to NATS server at %v by '%v'", f.serverAddr, conn.NatsConn().Opts.Name)
	return &Conn{
		conn:       conn,
		jStream:    js,
		serverAddr: f.serverAddr,
	}, nil
}
