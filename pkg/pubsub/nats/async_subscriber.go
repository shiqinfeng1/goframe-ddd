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

type AsyncSubscriber struct {
	logger pubsub.Logger
	nc     *nats.Conn
	subs   []*nats.Subscription
	ctx    context.Context
	cancel context.CancelFunc
}

func NewAsyncSubscriber(logger pubsub.Logger, f Factory) (*AsyncSubscriber, error) {
	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(gctx.New())

	// 连接NATS
	nc, err := f.New(ctx, "GoMgridAsyncSubscribeClient")
	if err != nil {
		cancel()
		return nil, err
	}

	return &AsyncSubscriber{
		logger: logger,
		nc:     nc,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (s *AsyncSubscriber) Subscribe(subject string, handler func(context.Context, *nats.Msg) error) error {
	sub, err := s.nc.Subscribe(subject, func(msg *nats.Msg) {
		err := func() error {
			defer recovery.Recovery(s.ctx, func(ctx context.Context, exception error) {
				s.logger.Errorf(ctx, "async subscriber: panic in handler: \n%v", exception)
			})
			if err := handler(s.ctx, msg); err != nil {
				time.Sleep(consumeMessageDelay)
				return err
			}
			return nil
		}()
		if err != nil {
			s.logger.Errorf(s.ctx, "async subscriber: error in handler for topic '%s': %v", subject, err)
		}
	})
	if err != nil {
		return err
	}
	sub.SetPendingLimits(-1, -1)
	s.subs = append(s.subs, sub)
	s.logger.Infof(s.ctx, "async subscribe for subject=%v ok", subject)
	return nil
}

func (s *AsyncSubscriber) Run() {
	s.logger.Infof(s.ctx, "async subscribe start running...")
	<-s.ctx.Done()
}

func (s *AsyncSubscriber) Shutdown() {
	s.logger.Infof(s.ctx, "async subscriber: starting shutdown...")

	// 取消所有订阅
	for _, sub := range s.subs {
		if err := sub.Drain(); err != nil {
			s.logger.Errorf(s.ctx, "async subscriber: error drain: %v", err)
		}
		if err := sub.Unsubscribe(); err != nil {
			s.logger.Errorf(s.ctx, "async subscriber: error unsubscribing: %v", err)
		}
	}
	s.nc.Close()
	s.logger.Infof(s.ctx, "async subscriber: shutdown complete")
	s.cancel()
}
