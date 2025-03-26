package pointmgr

import (
	"context"
	"time"

	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type PointDataSetMgr struct {
	repo Repository
}

// NewFileSendQueue 创建一个新的文件发送队列
func NewPointDataSetService(repo Repository) *PointDataSetMgr {
	return &PointDataSetMgr{
		repo: repo,
	}
}

func (p *PointDataSetMgr) HandleTopic1(ctx context.Context, msg *pubsub.Message) error {
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (p *PointDataSetMgr) HandleTopic2(ctx context.Context, msg *pubsub.Message) error {
	time.Sleep(20 * time.Millisecond)
	return nil
}
