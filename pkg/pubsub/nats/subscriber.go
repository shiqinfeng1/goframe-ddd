package nats

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

// 订阅器管理
type Subscriber struct {
	subscriptions map[string]*subscription
	subMutex      sync.Mutex
}

func NewSub() *Subscriber {
	sm := &Subscriber{
		subscriptions: make(map[string]*subscription),
	}
	return sm
}

func (sm *Subscriber) New(ctx context.Context, conn *nats.Conn, topicName string, consumeType SubType) error {
	sub := &subscription{
		conn:        conn,
		consumeType: consumeType,
		cancel:      make(chan struct{}),
		topicName:   topicName,
	}
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()
	if _, exist := sm.subscriptions[topicName]; exist {
		return gerror.Newf("topic '%v' is already be subscribed", topicName)
	}
	sm.subscriptions[topicName] = sub
	g.Log().Infof(ctx, "create subscriber of topic '%v' ok", topicName)
	return nil
}

func (sm *Subscriber) Close(ctx context.Context) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	for _, sub := range sm.subscriptions {
		if err := sub.unsubscribe(ctx); err != nil {
			return err
		}
	}
	sm.subscriptions = make(map[string]*subscription)
	return nil
}

func (sm *Subscriber) Delete(ctx context.Context, topicName string) error {
	sm.subMutex.Lock()

	sub, exist := sm.subscriptions[topicName]
	if !exist {
		sm.subMutex.Unlock()
		return gerror.Newf("not found subscription of topic '%v'", topicName)
	}
	sm.subscriptions = nil
	sm.subMutex.Unlock()

	if err := sub.unsubscribe(ctx); err != nil {
		return err
	}

	return nil
}
func (sm *Subscriber) Subscribe(ctx context.Context, topicName string, handler pubsub.SubscribeFunc) error {
	sm.subMutex.Lock()

	sub, exist := sm.subscriptions[topicName]
	if !exist {
		sm.subMutex.Unlock()
		return gerror.Newf("not found subscription of topic '%v'", topicName)
	}
	sm.subMutex.Unlock()

	if err := sub.Subscribe(ctx, handler); err != nil {
		return err
	}
	return nil
}

// func (sm *subscriber) createPubSubMessage(msg *nats.Msg, topic string) *nats.Msg {
// 	pubsubMsg := pubsub.NewMessage() // Pass a context if needed
// 	pubsubMsg.Topic = topic
// 	pubsubMsg.Value = msg.Data
// 	pubsubMsg.MetaData = msg.Header
// 	pubsubMsg.Subject = msg.Subject
// 	return pubsubMsg
// }
