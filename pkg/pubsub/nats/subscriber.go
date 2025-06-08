package natsclient

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

// 订阅器管理
type subscriber struct {
	logger        pubsub.Logger
	subscriptions *gmap.StrAnyMap //* subscription
}

func (sm *subscriber) AddSubscription(
	ctx context.Context,
	conn *Conn,
	topicName string,
	subTyp SubType,
	handler func(ctx context.Context, msg *nats.Msg) error) error {

	sub := NewSubscription(sm.logger, subTyp, topicName, conn, handler)

	if notexist := sm.subscriptions.SetIfNotExist(topicName, sub); !notexist {
		return gerror.Newf("topic '%v' is already be subscribed", topicName)
	}
	sm.logger.Infof(ctx, "create subscriber of topic '%v' ok", topicName)
	return nil
}

func (sm *subscriber) Close(ctx context.Context) error {
	if sm == nil {
		return nil
	}
	sm.subscriptions.Iterator(func(key string, value interface{}) bool {
		sub := value.(*subscription)
		if err := sub.Stop(ctx); err != nil {
			sm.logger.Errorf(ctx, "stop subscriber of topic '%v' failed: %v", key, err)
		}
		return true
	})

	sm.subscriptions.Clear()
	return nil
}

func (sm *subscriber) DeleteSubscription(ctx context.Context, topicName string) error {

	sub := sm.subscriptions.Remove(topicName)
	if sub == nil {
		return gerror.New("not found subscription of topic")
	}

	if err := sub.(*subscription).Stop(ctx); err != nil {
		return err
	}
	return nil
}
func (sm *subscriber) Start(ctx context.Context, topicName string) error {

	sub := sm.subscriptions.Get(topicName)
	if sub == nil {
		return gerror.Newf("not found subscription of topic '%v'", topicName)
	}

	if err := sub.(*subscription).Start(ctx); err != nil {
		return err
	}
	return nil
}
