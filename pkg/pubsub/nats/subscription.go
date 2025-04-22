package nats

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
	pubsub "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type subscription struct {
	consumeType SubType
	topicName   string
	conn        *nats.Conn
	sub         *nats.Subscription
	cancel      chan struct{}
}

func (s *subscription) unsubscribe(ctx context.Context) error {
	if s.sub != nil {
		if err := s.sub.Drain(); err != nil {
			return err
		}
		if err := s.sub.Unsubscribe(); err != nil {
			return err
		}
	}
	close(s.cancel)
	g.Log().Infof(ctx, "unsubscribe topic '%v' ok", s.topicName)
	return nil
}
func (s *subscription) Subscribe(ctx context.Context, handler pubsub.SubscribeFunc) error {
	switch s.consumeType {
	case SubTypeSubAsync:
		return s.subscribeAsync(ctx, handler)
	}
	return nil
}

// 订阅指定的topic的消息
func (s *subscription) subscribeAsync(
	ctx context.Context,
	handler pubsub.SubscribeFunc,
) error {
	sub, err := s.conn.Subscribe(s.topicName, func(msg *nats.Msg) {
		metrics.IncrementCounter(ctx, metrics.NatsSubscribeTotalCount, "topic", s.topicName)

		err := func() error {
			defer func() {
				panicRecovery(ctx, recover())
			}()
			if err := handler(ctx, msg); err != nil {
				time.Sleep(ConsumeMessageDelay)
				return err
			}
			return nil
		}()
		if err != nil {
			g.Log().Errorf(ctx, "error in handler for topic '%s': %v", s.topicName, err)
		}

	})
	if err != nil {
		return err
	}
	sub.SetPendingLimits(-1, -1)
	s.sub = sub
	g.Log().Infof(ctx, "ready to subscribe msg for '%v'", s.topicName)
	// 等待被取消
	<-s.cancel
	g.Log().Infof(ctx, "subscribe of '%v' is canceled", s.topicName)
	return nil
}
