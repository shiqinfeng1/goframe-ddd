package natsclient

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
	"github.com/shiqinfeng1/goframe-ddd/pkg/recovery"
)

// 订阅器管理

type SyncSubscriber struct {
	logger  pubsub.Logger
	nc      *nats.Conn
	subs    []*nats.Subscription
	ctx     context.Context
	cancel  context.CancelFunc
	errChan chan error
}

func NewSyncSubscriber(logger pubsub.Logger, f Factory) (*SyncSubscriber, error) {
	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(gctx.New())

	// 连接NATS
	nc, err := f.New(ctx, "GoMgridSyncSubscribeClient")
	if err != nil {
		cancel()
		return nil, err
	}

	return &SyncSubscriber{
		logger:  logger,
		nc:      nc,
		ctx:     ctx,
		cancel:  cancel,
		errChan: make(chan error, 100),
	}, nil
}

func (s *SyncSubscriber) Subscribe(subject string, handler func(ctx context.Context, msg *nats.Msg) error) error {
	sub, err := s.nc.SubscribeSync(subject)
	if err != nil {
		return err
	}
	sub.SetPendingLimits(-1, -1)
	s.subs = append(s.subs, sub)
	go func() {
		s.processMessages(sub, handler)
	}()
	s.logger.Infof(s.ctx, "sync subscribe for subject=%v ok", subject)
	return nil
}
func (s *SyncSubscriber) processMessages(sub *nats.Subscription, handler func(ctx context.Context, msg *nats.Msg) error) {
	timeout := 30 * time.Second
	for {
		msg, err := sub.NextMsg(timeout)
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			}
			s.errChan <- err
			return
		}
		err = func() error {
			defer recovery.Recovery(s.ctx, func(ctx context.Context, exception error) {
				s.logger.Errorf(ctx, "sync subscriber: panic in nats handler: \n%v", exception)
			})
			if err := handler(s.ctx, msg); err != nil {
				time.Sleep(consumeMessageDelay)
				return err
			}
			return nil
		}()
		if err != nil {
			s.logger.Errorf(s.ctx, "sync subscriber: error in handler for topic '%s': %v", sub.Subject, err)
		}
	}
}
func (s *SyncSubscriber) Run() {
	go func() {
		for err := range s.errChan {
			s.logger.Errorf(s.ctx, "sync subscriber error: %v", err)
		}
	}()
	s.logger.Infof(s.ctx, "sync subscribe start running...")
	<-s.ctx.Done()
}

func (s *SyncSubscriber) Shutdown() {
	s.logger.Infof(s.ctx, "async subscriber: starting shutdown...")

	// 取消上下
	s.cancel()

	// 取消所有订阅
	for _, sub := range s.subs {
		if err := sub.Drain(); err != nil {
			s.logger.Infof(s.ctx, "sync subscriber: error drain: %v", err)
		}
		if err := sub.Unsubscribe(); err != nil {
			s.logger.Infof(s.ctx, "sync subscriber: error unsubscribing: %v", err)
		}
	}
	// 关闭错误通道
	close(s.errChan)
	// 关闭连接
	s.nc.Close()
	s.logger.Infof(s.ctx, "async subscriber: shutdown complete")
}
