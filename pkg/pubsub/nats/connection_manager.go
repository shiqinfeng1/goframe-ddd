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

//go:generate mockgen -destination=mock_jetstream.go -package=nats github.com/nats-io/nats.go/jetstream jStream,Stream,Consumer,Msg,MessageBatch

type ConnectionManager struct {
	conn             ConnIntf
	jStream          jetstream.JetStream
	config           *Config
	natsConnector    Connector
	jetStreamCreator JetStreamCreator
}

func (cm *ConnectionManager) jetStream() (jetstream.JetStream, error) {
	if cm.jStream == nil {
		return nil, errJetStreamNotConfigured
	}

	return cm.jStream, nil
}

// NewConnectionManager creates a new ConnectionManager.
func newConnectionManager(
	cfg *Config,
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
		config:           cfg,
		natsConnector:    natsConnector,
		jetStreamCreator: jetStreamCreator,
	}
}

// Connect establishes a connection to NATS and sets up JetStream.
// 异步重试连接
func (cm *ConnectionManager) Connect(ctx context.Context) error {
	go cm.retryConnect(ctx)
	return nil
}

func (cm *ConnectionManager) retryConnect(ctx context.Context) {
	opts := []nats.Option{nats.Name("Sieyuan NATS JetStreamClient")}

	for {
		connIntf, err := cm.natsConnector.Connect(cm.config.Server, opts...)
		if err != nil {
			g.Log().Errorf(ctx, "Failed to connect to NATS server at %v: %v", cm.config.Server, err)
			time.Sleep(defaultRetryTimeout) // 等待10s后再连
			continue
		}
		// 连接成功后，创建jetstream实例
		js, err := cm.jetStreamCreator.New(connIntf)
		if err != nil {
			connIntf.Close()
			g.Log().Debugf(ctx, "Failed to create jStream context: %v", err)
			time.Sleep(defaultRetryTimeout)
			continue
		}

		cm.conn = connIntf
		cm.jStream = js
		g.Log().Infof(ctx, "Successfully connected to NATS server at %v", cm.config.Server)
		return
	}
}

func (cm *ConnectionManager) Close(_ context.Context) {
	if cm.conn != nil {
		cm.conn.Close()
	}
}

func (cm *ConnectionManager) Publish(ctx context.Context, subject string, message []byte) error {
	metrics.IncrementCounter(ctx, metrics.NatsPublishTotalCount, "subject", subject)

	if err := cm.validateJetStream(ctx, subject); err != nil {
		return err
	}
	// 发布消息
	_, err := cm.jStream.Publish(ctx, subject, message)
	if err != nil {
		g.Log().Errorf(ctx, "failed to publish message to NATS jStream: %v", err)
		return err
	}

	metrics.IncrementCounter(ctx, metrics.NatsPublishSuccessCount, "subject", subject)

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
func (cm *ConnectionManager) Health() health.Health {
	if cm.conn == nil {
		return health.Health{
			Status: health.StatusDown,
		}
	}

	status := cm.conn.Status()
	if status == nats.CONNECTED {
		return health.Health{
			Status: health.StatusUp,
			Details: map[string]any{
				"server": cm.config.Server,
			},
		}
	}

	return health.Health{
		Status: health.StatusDown,
		Details: map[string]any{
			"server": cm.config.Server,
		},
	}
}

func (cm *ConnectionManager) isConnected() bool {
	if cm.conn == nil {
		return false
	}

	return cm.conn.Status() == nats.CONNECTED
}
