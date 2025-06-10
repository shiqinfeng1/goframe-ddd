package natsclient

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
	"github.com/shiqinfeng1/goframe-ddd/pkg/recover"
)

type closer struct {
	cancel context.CancelFunc
	drain  func()
}
type streamConsume struct {
	logger      pubsub.Logger
	consumeType SubType
	subsKey     SubsKey
	consumer    jetstream.Consumer
	handler     func(ctx context.Context, msg *jetstream.Msg) error
	close       closer
	exitNotify  chan SubsKey
}

func NewStreamConsume(
	l pubsub.Logger,
	st SubType,
	sk SubsKey,
	c jetstream.Consumer,
	handler func(ctx context.Context, msg *jetstream.Msg) error,
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
	if s.close.drain != nil {
		s.close.drain()
	}
	if s.close.cancel != nil {
		s.close.cancel()
	}
	return nil
}

func (s *streamConsume) start(ctx context.Context) error {
	switch s.consumeType {
	case JSNEXT:
		return s.consumeNext(ctx)
	case JSFETCH:
	}
	return nil
}

func (s *streamConsume) consumeNext(
	ctx context.Context,
) error {
	metrics.Inc(ctx, metrics.NatsSubscribeTotalCount, "topic", s.subsKey.TopicName())

	iter, err := s.newMessageIter()
	if err != nil {
		return err
	}
	defer func() {
		s.exitNotify <- s.subsKey
	}()
	s.logger.Infof(ctx, "[consumeNext]start consume. stream-name=%v, consumer=%v, subject=%v", s.subsKey.StreamName(), s.subsKey.ConsumerName(), s.subsKey.TopicName())
	// 获取主题对应的消息队列缓存
	for {
		msg, err := iter.Next()
		if err != nil {
			if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
				s.logger.Warningf(ctx, "[consumeNext]msg iter closed. stream-name=%v, consumer=%v, subject=%v", s.subsKey.StreamName(), s.subsKey.ConsumerName(), msg.Subject())
				return nil
			}
			s.logger.Warningf(ctx, "[consumeNext]get next msg fail. stream-name=%v, consumer=%v, subject=%v: %v", s.subsKey.StreamName(), s.subsKey.ConsumerName(), msg.Subject(), err)
			time.Sleep(consumeMessageDelay)
			iter, err = s.newMessageIter()
			if err != nil {
				return err
			}
			continue
		}
		err = func() error {
			defer func() {
				recover.Recovery(ctx, func(ctx context.Context, exception error) {
					s.logger.Errorf(ctx, "panic in handler:%v", exception)
				})
			}()
			if err := s.handler(ctx, &msg); err != nil {
				time.Sleep(consumeMessageDelay)
				return err
			}
			return nil
		}()
		if err != nil {
			s.logger.Errorf(ctx, "[consumeNext]handler fail. stream-name=%v, consumer=%v, subject=%v: %v", s.subsKey.StreamName(), s.subsKey.ConsumerName(), msg.Subject(), err)
			continue
		}
		// 处理完成
		if err := msg.Ack(); err != nil {
			s.logger.Errorf(ctx, "[consumeNext]ack fail. stream-name=%v, consumer=%v, subject=%v: %v", s.subsKey.StreamName(), s.subsKey.ConsumerName(), msg.Subject(), err)
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
		drain: iter.Drain, // drain 保证本地缓存中的消息被处理完之后才会关闭iter
	}
	return iter, nil
}
