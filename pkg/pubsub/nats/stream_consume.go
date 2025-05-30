package nats

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
	"github.com/shiqinfeng1/goframe-ddd/pkg/panic"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type closer struct {
	cancel context.CancelFunc
	stop   func()
}
type streamConsume struct {
	logger      pubsub.Logger
	consumeType SubType
	subsKey     SubsKey
	consumer    jetstream.Consumer
	handler     ConsumeFunc
	close       closer
	exitNotify  chan SubsKey
}

func NewStreamConsume(
	l pubsub.Logger,
	st SubType,
	sk SubsKey,
	c jetstream.Consumer,
	handler ConsumeFunc,
	exit chan SubsKey) *streamConsume {

	return &streamConsume{
		logger:      l,
		consumeType: st,
		subsKey:     sk,
		consumer:    c,
		handler:     handler,
		exitNotify:  exit,
	}
}

func (s *streamConsume) Stop(ctx context.Context) error {
	if s.close.stop != nil {
		s.close.stop()
	}
	if s.close.cancel != nil {
		if s.close.cancel != nil {
			s.close.cancel()
		}
	}
	return nil
}

func (s *streamConsume) start(ctx context.Context) error {
	switch s.consumeType {
	case SubTypeJSConsumeNext:
		return s.consumeNext(ctx)
	}
	return nil
}

func (s *streamConsume) consumeNext(
	ctx context.Context,
) error {
	metrics.IncCnt(ctx, metrics.NatsSubscribeTotalCount, "topic", s.subsKey.TopicName())

	iter, err := s.newMessageIter()
	if err != nil {
		return err
	}
	defer func() {
		s.exitNotify <- s.subsKey
	}()
	s.logger.Infof(ctx, "consumer '%v' ready to consume msg for '%v' of '%v'", s.subsKey.ConsumerName(), s.subsKey.TopicName(), s.subsKey.StreamName())
	// 获取主题对应的消息队列缓存
	for {
		msg, err := iter.Next()
		if err != nil {
			if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
				s.logger.Warningf(ctx, "consumer '%v' fetching messages for topic '%s' of '%v': %v", s.subsKey.ConsumerName(), s.subsKey.TopicName(), s.subsKey.StreamName(), err)
				iter, err = s.newMessageIter()
				if err != nil {
					return err
				}
				time.Sleep(ConsumeMessageDelay)
				s.logger.Warningf(ctx, "consumer '%v' subscribe messages again for topic '%s' of '%v' ok", s.subsKey.ConsumerName(), s.subsKey.TopicName(), s.subsKey.StreamName())
				continue
			}
			s.logger.Warningf(ctx, "consumer '%v' fetching messages for topic '%s' of '%v' fail: %v", s.subsKey.ConsumerName(), s.subsKey.TopicName(), s.subsKey.StreamName(), err)
			return nil
		}
		err = func() error {
			defer func() {
				panic.Recovery(ctx, func(ctx context.Context, exception error) {
					s.logger.Errorf(ctx, "panic in handler:%v", exception)
				})
			}()
			if err := s.handler(ctx, &msg); err != nil {
				time.Sleep(ConsumeMessageDelay)
				return err
			}
			return nil
		}()
		if err != nil {
			s.logger.Errorf(ctx, "consumer '%v' error in handler for subject '%s': %v", s.subsKey.ConsumerName(), msg.Subject(), err)
			continue
		}
		// 处理完成
		if err := msg.Ack(); err != nil {
			s.logger.Errorf(ctx, "consumer '%v' ack fail for subject '%s': %v", s.subsKey.ConsumerName(), msg.Subject(), err)
		}
	}
}

func (s *streamConsume) newMessageIter() (jetstream.MessagesContext, error) {
	// 获取消息迭代器
	iter, err := s.consumer.Messages(jetstream.PullMaxMessages(1))
	if err != nil {
		return nil, gerror.Wrap(err, "get consumer msg iter fail")
	}
	// 注册订阅者
	s.close = closer{
		stop: iter.Stop,
	}
	return iter, err
}
