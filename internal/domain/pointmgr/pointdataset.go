package pointmgr

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type PointDataSetMgr struct {
	repo Repository
}

// NewFileSendQueue 创建一个新的文件发送队列
func NewPointDataSetService(_ context.Context, repo Repository) *PointDataSetMgr {
	return &PointDataSetMgr{
		repo: repo,
	}
}

func (p *PointDataSetMgr) HandleTopic1(ctx context.Context, msg *nats.Msg) error {
	time.Sleep(10 * time.Millisecond)
	// g.Log().Debugf(ctx, "1 recv a msg: %v", msg)
	return nil
}

func (p *PointDataSetMgr) HandleTopic2(ctx context.Context, msg *jetstream.Msg) error {
	time.Sleep(5 * time.Millisecond)
	// g.Log().Debugf(ctx, "2 recv a msg: %v", msg)
	return nil
}
