package service

import (
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/repository"
)

type PointDataSetMgr struct {
	repo   repository.Repository
	logger application.Logger
}

// NewFileSendQueue 创建一个新的文件发送队列
func NewPointDataSetService(_ context.Context, logger application.Logger, repo repository.Repository) application.PointDataSetSrv {
	return &PointDataSetMgr{
		logger: logger,
		repo:   repo,
	}
}

func (p *PointDataSetMgr) HandleMsg(ctx context.Context, msg *nats.Msg) ([]byte, error) {
	time.Sleep(10 * time.Millisecond)
	// g.Log().Debugf(ctx, "1 recv a msg: %v", msg)
	return (*msg).Data, nil
}

func (p *PointDataSetMgr) HandleStream(ctx context.Context, msg *jetstream.Msg) ([]byte, error) {
	time.Sleep(5 * time.Millisecond)
	// g.Log().Debugf(ctx, "2 recv a msg: %v", msg)
	return (*msg).Data(), nil
}
func (p *PointDataSetMgr) HandleMqttMsg(ctx context.Context, msg *mqtt.Message) ([]byte, error) {
	time.Sleep(5 * time.Millisecond)
	// g.Log().Debugf(ctx, "2 recv a msg: %v", msg)
	return (*msg).Payload(), nil
}
