package nats

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
)

//go:generate mockgen -destination=mock_jetstream.go -package=nats github.com/nats-io/nats.go/jetstream JetStream,Stream,Consumer,Msg,MessageBatch

type ConnectionManager struct {
	conn             ConnIntf
	jStream          jetstream.JetStream
	serverAddr       string // 服务地址
	natsConnector    Connector
	jetStreamCreator JetStreamCreator
}

func (cm *ConnectionManager) GetJetStream() (jetstream.JetStream, error) {
	if cm.jStream == nil {
		return nil, errJetStreamNotConfigured
	}

	return cm.jStream, nil
}

// newConnMgr creates a new ConnectionManager.
func newConnMgr(
	serverAddr string,
	natsConnector Connector,
	jetStreamCreator JetStreamCreator,
) *ConnectionManager {
	// 设置连接器
	if natsConnector == nil {
		natsConnector = &defaultConnector{}
	}
	// 设置js构造器
	if jetStreamCreator == nil {
		jetStreamCreator = &defaultJetStreamCreator{}
	}

	return &ConnectionManager{
		natsConnector:    natsConnector,
		jetStreamCreator: jetStreamCreator,
		serverAddr:       serverAddr,
	}
}

// Connect establishes a connection to NATS and sets up JetStream.
// 异步重试连接
func (cm *ConnectionManager) Connect(ctx context.Context, opts ...nats.Option) {
	for {
		conn, err := cm.natsConnector.Connect(cm.serverAddr, opts...)
		if err != nil {
			g.Log().Warningf(ctx, "try to connect to NATS server at %v: %v", cm.serverAddr, err)
			time.Sleep(defaultRetryTimeout)
			continue
		}
		// 连接成功后，创建jetstream
		js, err := cm.jetStreamCreator.New(conn)
		if err != nil {
			conn.Close()
			g.Log().Debugf(ctx, "Failed to create jStream context: %v", err)
			time.Sleep(defaultRetryTimeout)
			continue
		}

		cm.conn = conn
		cm.jStream = js
		g.Log().Infof(ctx, "Successfully connected to NATS server at %v by '%v'", cm.serverAddr, conn.Conn().Opts.Name)
		return
	}
}

func (cm *ConnectionManager) Close(_ context.Context) {
	if cm.conn != nil {
		cm.conn.Close()
	}
}

func (cm *ConnectionManager) Subscribe(ctx context.Context, subject string, handler func(msg *nats.Msg)) error {
	metrics.IncrementCounter(ctx, metrics.NatsPublishTotalCount, "subject", subject)
	subs, err := cm.conn.Conn().Subscribe(subject, handler)
	if err != nil {
		return err
	}
	subs.Unsubscribe()
	metrics.IncrementCounter(ctx, metrics.NatsPublishSuccessCount, "subject", subject)
	return nil
}
func (cm *ConnectionManager) Publish(ctx context.Context, subject string, message []byte) error {
	metrics.IncrementCounter(ctx, metrics.NatsPublishTotalCount, "subject", subject)
	if err := cm.conn.Conn().Publish(subject, message); err != nil {
		return err
	}
	metrics.IncrementCounter(ctx, metrics.NatsPublishSuccessCount, "subject", subject)
	return nil
}
func (cm *ConnectionManager) JsPublish(ctx context.Context, subject string, message []byte) error {
	metrics.IncrementCounter(ctx, metrics.NatsJsPublishTotalCount, "subject", subject)

	if err := cm.validateJetStream(ctx, subject); err != nil {
		return err
	}
	// 发布消息
	_, err := cm.jStream.Publish(ctx, subject, message)
	if err != nil {
		g.Log().Errorf(ctx, "failed to publish message to NATS jStream: %v", err)
		return err
	}

	metrics.IncrementCounter(ctx, metrics.NatsJsPublishSuccessCount, "subject", subject)

	return nil
}

func (cm *ConnectionManager) validateJetStream(_ context.Context, subject string) error {
	if cm.jStream == nil || subject == "" {
		err := errJetStreamNotConfigured
		return err
	}

	return nil
}

// 返回nats客户端和服务端之间连接的健康状态
func (cm *ConnectionManager) Health() *health.Health {
	if cm.conn == nil {
		return &health.Health{
			Status: health.StatusDown,
		}
	}

	status := cm.conn.Status()
	if status == nats.CONNECTED {
		return &health.Health{
			Status: health.StatusUp,
			Details: map[string]any{
				"server": cm.serverAddr,
			},
		}
	}

	return &health.Health{
		Status: health.StatusDown,
		Details: map[string]any{
			"server": cm.serverAddr,
		},
	}
}

func (cm *ConnectionManager) isConnected() bool {
	if cm.conn == nil {
		return false
	}

	return cm.conn.Status() == nats.CONNECTED
}
