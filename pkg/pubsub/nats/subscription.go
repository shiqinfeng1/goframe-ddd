package natsclient

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
	"github.com/shiqinfeng1/goframe-ddd/pkg/recover"
)

type subscription struct {
	logger    pubsub.Logger
	stype     SubType
	topicName string
	conn      *Conn
	sub       *nats.Subscription
	handler   func(ctx context.Context, msg *nats.Msg) error
	exit      chan struct{}
}

func NewSubscription(
	l pubsub.Logger,
	stype SubType,
	tn string,
	c *Conn,
	handler func(ctx context.Context, msg *nats.Msg) error) *subscription {

	return &subscription{
		logger:    l,
		stype:     stype,
		topicName: tn,
		conn:      c,
		handler:   handler,
		exit:      make(chan struct{}),
	}
}
func (s *subscription) Stop(ctx context.Context) error {
	if s.sub != nil {
		if err := s.sub.Drain(); err != nil {
			return gerror.Wrap(err, "nats drain fail")
		}
		if err := s.sub.Unsubscribe(); err != nil {
			return gerror.Wrap(err, "nats unsub fail")
		}
	}
	close(s.exit)
	s.logger.Infof(ctx, "unsubscribe topic '%v' ok", s.topicName)
	return nil
}
func (s *subscription) Start(ctx context.Context) error {
	switch s.stype {
	case SUBASYNC:
		return s.subscribeAsync(ctx, s.handler)
	case SUBSYNC:
	}
	return nil
}

// 订阅指定的topic的消息
func (s *subscription) subscribeAsync(
	ctx context.Context,
	handler func(ctx context.Context, msg *nats.Msg) error,
) error {
	sub, err := s.conn.SubMsg(ctx, s.topicName, func(msg *nats.Msg) {
		metrics.Inc(ctx, metrics.NatsSubscribeTotalCount, "topic", s.topicName)

		err := func() error {
			defer func() {
				recover.Recovery(ctx, func(ctx context.Context, exception error) {
					s.logger.Errorf(ctx, "panic in nats handler:%v", exception)
				})
			}()
			if err := handler(ctx, msg); err != nil {
				time.Sleep(consumeMessageDelay)
				return err
			}
			return nil
		}()
		if err != nil {
			s.logger.Errorf(ctx, "error in handler for topic '%s': %v", s.topicName, err)
		}
	})
	if err != nil {
		return err
	}
	sub.SetPendingLimits(-1, -1)
	s.sub = sub
	s.logger.Infof(ctx, "ready to subscribe msg for '%v'", s.topicName)
	// 等待被取消
	<-s.exit
	s.logger.Infof(ctx, "subscribe of '%v' is exited", s.topicName)
	return nil
}
