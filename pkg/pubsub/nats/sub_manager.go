package nats

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type subscriber struct {
	consumeType SubType
	topicName   string
	conn        *nats.Conn
	sub         *nats.Subscription
	cancel      chan struct{}
}

// 订阅器管理
type SubscriberManager struct {
	subscriptions map[string]*subscriber
	subMutex      sync.Mutex
}

func NewSubMgr() *SubscriberManager {
	sm := &SubscriberManager{
		subscriptions: make(map[string]*subscriber),
	}
	return sm
}

func (sm *SubscriberManager) NewSubscriber(ctx context.Context, conn *nats.Conn, topicName string, consumeType SubType) error {
	sub := &subscriber{
		conn:        conn,
		consumeType: consumeType,
		cancel:      make(chan struct{}),
	}
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()
	if _, exist := sm.subscriptions[topicName]; exist {
		return gerror.Newf("topic '%v' is already be subscribed", topicName)
	}
	sm.subscriptions[topicName] = sub
	g.Log().Infof(ctx, "create subscriber ok. topicName: '%v'", topicName)
	return nil
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
	if s.sub != nil {
		if err := s.sub.Drain(); err != nil {
			return err
		}
		if err := s.sub.Unsubscribe(); err != nil {
			return err
		}
		return nil
	}
	close(s.cancel)
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
	sub, err := s.conn.Subscribe(s.topicName, func(msg *nats.Msg) {
		metrics.IncrementCounter(ctx, metrics.NatsSubscribeTotalCount, "topic", s.topicName)

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
	sub.SetPendingLimits(-1, -1)
	s.sub = sub
	// 等待被取消
	<-s.cancel
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
