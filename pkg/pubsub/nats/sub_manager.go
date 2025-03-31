package nats

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type subscription struct {
	stop func() error
}
type subscriber struct {
	consumeType SubType
	topicName   string
	conn        *nats.Conn
	closer      subscription
}

// 订阅器管理
type SubscriberManager struct {
	subscriptions map[string]*subscriber
	subMutex      sync.Mutex
}

func NewSubscriberManager() *SubscriberManager {
	sm := &SubscriberManager{
		subscriptions: make(map[string]*subscriber),
	}
	return sm
}

func (sm *SubscriberManager) NewSubscriber(conn *nats.Conn, topicName string, consumeType SubType) *subscriber {
	sub := &subscriber{
		conn:        conn,
		consumeType: consumeType,
	}
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()
	sm.subscriptions[topicName] = sub
	return sub
}

func (sm *SubscriberManager) Close(ctx context.Context) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	for _, sub := range sm.subscriptions {
		if err := sub.unsubscribe(ctx); err != nil {
			return err
		}
	}
	sm.subscriptions = make(map[string]*subscriber)
	return nil
}
func (sm *SubscriberManager) deleteSubscriber(ctx context.Context, key string) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	if sub, exist := sm.subscriptions[key]; exist {
		if err := sub.unsubscribe(ctx); err != nil {
			return err
		}
		sm.subscriptions = nil
	}
	return nil
}
func (s *subscriber) unsubscribe(ctx context.Context) error {
	if s.closer.stop != nil {
		s.closer.stop()
	}
	g.Log().Infof(ctx, "unsubscribe topic <%v> ok", s.topicName)
	return nil
}
func (sm *SubscriberManager) DeleteSubscriber(ctx context.Context, topicName string) error {
	return sm.deleteSubscriber(ctx, topicName)
}
func (sm *SubscriberManager) Subscribe(ctx context.Context, topicName string, handler pubsub.SubscribeFunc) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	if sub, exist := sm.subscriptions[topicName]; exist {
		if err := sub.Subscribe(ctx, handler); err != nil {
			return err
		}
	}
	return nil
}

func (s *subscriber) Subscribe(ctx context.Context, handler pubsub.SubscribeFunc) error {
	switch s.consumeType {
	case SubTypeSubAsync:
		return s.subscribeAsync(ctx, handler)
	}
	return nil
}

// 订阅指定的topic的消息
func (s *subscriber) subscribeAsync(
	ctx context.Context,
	handler pubsub.SubscribeFunc,
) error {
	metrics.IncrementCounter(ctx, metrics.NatsSubscribeTotalCount, "topic", s.topicName)

	subs, err := s.conn.Subscribe(s.topicName, func(msg *nats.Msg) {
		pubsubmsg := s.createPubSubMessage(msg, s.topicName)
		err := func() error {
			defer func() {
				panicRecovery(ctx, recover())
			}()
			if err := handler(ctx, pubsubmsg); err != nil {
				time.Sleep(ConsumeMessageDelay)
				return err
			}
			return nil
		}()
		if err != nil {
			g.Log().Errorf(ctx, "error in handler for topic '%s': %v", s.topicName, err)
		}
		// 处理完成
		if pubsubmsg.Committer != nil {
			pubsubmsg.Commit(ctx)
		}
	})
	if err != nil {
		return err
	}
	s.closer = subscription{stop: subs.Unsubscribe}
	return nil
}

func (sm *subscriber) createPubSubMessage(msg *nats.Msg, topic string) *pubsub.Message {
	pubsubMsg := pubsub.NewMessage() // Pass a context if needed
	pubsubMsg.Topic = topic
	pubsubMsg.Value = msg.Data
	pubsubMsg.MetaData = msg.Header
	pubsubMsg.Committer = &natsCommitter{msg: msg}
	return pubsubMsg
}
