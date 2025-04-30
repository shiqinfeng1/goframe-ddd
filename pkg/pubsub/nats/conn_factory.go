package nats

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type ConnFactory interface {
	New(ctx context.Context, opts ...nats.Option) (*Conn, error)
}
type factory struct {
	logger     pubsub.Logger
	serverAddr string
	connector  Connector
}

func NewFactory(
	logger pubsub.Logger,
	serverAddr string,
	natsConnector Connector,
) ConnFactory {
	// 设置连接器
	if natsConnector == nil {
		natsConnector = &defaultConnector{}
	}

	return &factory{
		logger:     logger,
		connector:  natsConnector,
		serverAddr: serverAddr,
	}
}
func (f *factory) New(ctx context.Context, opts ...nats.Option) (*Conn, error) {
	for i := range defaultRetryCount {
		conn, err := f.connector.Connect(f.serverAddr, opts...)
		if err != nil {
			f.logger.Warningf(ctx, "[%v/%v]try to connect to NATS server at %v: %v", i+1, defaultRetryCount, f.serverAddr, err)
			time.Sleep(defaultRetryTimeout)
			continue
		}
		// 连接成功后，创建jetstream
		js, err := conn.NewJetStream()
		if err != nil {
			conn.Close()
			f.logger.Debugf(ctx, "[%v/%v]Failed to create jStream context: %v", i+1, defaultRetryCount, err)
			time.Sleep(defaultRetryTimeout)
			continue
		}
		f.logger.Infof(ctx, "Successfully connected to NATS server at %v by '%v'", f.serverAddr, conn.NatsConn().Opts.Name)
		return &Conn{
			conn:       conn,
			jStream:    js,
			serverAddr: f.serverAddr,
		}, nil
	}
	return nil, gerror.New("connect to nats timeout")
}
