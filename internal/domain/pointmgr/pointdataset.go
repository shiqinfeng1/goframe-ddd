package pointmgr

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
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
	time.Sleep(5 * time.Millisecond)
	g.Log().Debugf(ctx, "1 recv a msg: %v", msg)
	return nil
}

func (p *PointDataSetMgr) HandleTopic2(ctx context.Context, msg *pubsub.Message) error {
	time.Sleep(5 * time.Millisecond)
	g.Log().Debugf(ctx, "2 recv a msg: %v", msg)
	return nil
}
