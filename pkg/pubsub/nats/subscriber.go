package nats

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

// 订阅器管理
type Subscriber struct {
	logger        pubsub.Logger
	subscriptions map[string]*subscription
	subMutex      sync.Mutex
}

func NewSubscriber(logger pubsub.Logger) *Subscriber {
	sm := &Subscriber{
		logger:        logger,
		subscriptions: make(map[string]*subscription),
		subMutex:      sync.Mutex{},
	}
	return sm
}

func (sm *Subscriber) AddSubscription(
	ctx context.Context,
	conn *Conn,
	topicName string,
	subTyp SubType,
	handler SubscribeFunc) error {

	sub := NewSubscription(sm.logger, subTyp, topicName, conn, handler)

	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()
	if _, exist := sm.subscriptions[topicName]; exist {
		return gerror.Newf("topic '%v' is already be subscribed", topicName)
	}
	sm.subscriptions[topicName] = sub
	sm.logger.Infof(ctx, "create subscriber of topic '%v' ok", topicName)
	return nil
}

func (sm *Subscriber) Close(ctx context.Context) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()
	for _, sub := range sm.subscriptions {
		if err := sub.Stop(ctx); err != nil {
			return err
		}
	}
	sm.subscriptions = make(map[string]*subscription)
	return nil
}

func (sm *Subscriber) DeleteSub(ctx context.Context, topicName string) error {

	sm.subMutex.Lock()
	sub, exist := sm.subscriptions[topicName]
	if !exist {
		sm.subMutex.Unlock()
		return gerror.Newf("not found subscription of topic '%v'", topicName)
	}
	sm.subscriptions = nil
	sm.subMutex.Unlock()

	if err := sub.Stop(ctx); err != nil {
		return err
	}
	return nil
}
func (sm *Subscriber) StartSub(ctx context.Context, topicName string) error {

	sm.subMutex.Lock()
	sub, exist := sm.subscriptions[topicName]
	if !exist {
		sm.subMutex.Unlock()
		return gerror.Newf("not found subscription of topic '%v'", topicName)
	}
	sm.subMutex.Unlock()

	if err := sub.Start(ctx); err != nil {
		return err
	}
	return nil
}
