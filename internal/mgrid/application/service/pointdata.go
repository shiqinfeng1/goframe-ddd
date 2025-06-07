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

type PointdataService struct {
	repo   repository.PointdataRepository
	logger application.Logger
}

// NewFileSendQueue 创建一个新的文件发送队列
func NewPointdataService(_ context.Context, logger application.Logger, repo repository.PointdataRepository) application.PointdataService {
	return &PointdataService{
		logger: logger,
		repo:   repo,
	}
}

func (p *PointdataService) HandleMsg(ctx context.Context, msg *nats.Msg) ([]byte, error) {
	time.Sleep(10 * time.Millisecond)
	p.logger.Debugf(ctx, "recv a msg: %+v", msg)
	return (*msg).Data, nil
}

func (p *PointdataService) HandleStream(ctx context.Context, msg *jetstream.Msg) ([]byte, error) {
	time.Sleep(5 * time.Millisecond)
	p.logger.Debugf(ctx, "recv a stream data: %+v", msg)
	return (*msg).Data(), nil
}
func (p *PointdataService) HandleMqttMsg(ctx context.Context, msg *mqtt.Message) ([]byte, error) {
	time.Sleep(5 * time.Millisecond)
	p.logger.Debugf(ctx, "recv a mqtt msg: %+v", msg)
	return (*msg).Payload(), nil
}
