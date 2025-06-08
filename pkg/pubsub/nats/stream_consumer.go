package natsclient

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

// 订阅器管理
type consumer struct {
	logger        pubsub.Logger
	subscriptions *gmap.StrAnyMap //map[string]*streamConsume
	streamMgr     *JetStreamWrapper
	exitNotify    chan SubsKey // 同步模式下。订阅退出后，各个consume通知本consumer删除相关记录
}

func (sm *consumer) Add(
	ctx context.Context,
	st SubType,
	sk SubsKey,
	c jetstream.Consumer,
	handler func(ctx context.Context, msg *jetstream.Msg) error,
	exit chan SubsKey) error {

	ssub := NewStreamConsume(sm.logger, st, sk, c, handler, exit)

	if notexist := sm.subscriptions.SetIfNotExist(sk.String(), ssub); !notexist {
		return gerror.Newf("stream '%v' topic '%v' is already be consumed by '%v'", sk.StreamName(), sk.TopicName(), sk.ConsumerName())
	}
	if err := ssub.start(ctx); err != nil {
		return err
	}
	sm.logger.Infof(ctx, "stream '%v' create consumer '%v' of topic '%v' ok", sk.StreamName(), sk.ConsumerName(), sk.TopicName())
	return nil
}

func (sm *consumer) Close(ctx context.Context) error {
	if sm == nil {
		return nil
	}
	sm.subscriptions.Iterator(func(key string, value interface{}) bool {
		sub := value.(*streamConsume)
		if err := sm.streamMgr.DeleteConsumer(ctx, sub.subsKey.StreamName(), sub.subsKey.ConsumerName(), sub.subsKey.TopicName()); err != nil {
			sm.logger.Errorf(ctx, "delete cunsumer <%v> for topic <%v> of stream <%v> failed: %v", sub.subsKey.ConsumerName(), sub.subsKey.TopicName(), sub.subsKey.StreamName(), err)
		}
		if err := sub.Stop(ctx); err != nil {
			sm.logger.Errorf(ctx, "stop cunsumer <%v> for topic <%v> of stream <%v> failed: %v", sub.subsKey.ConsumerName(), sub.subsKey.TopicName(), sub.subsKey.StreamName(), err)
		}
		return true
	})

	sm.subscriptions.Clear()
	close(sm.exitNotify)
	return nil
}

func (sm *consumer) Delete(ctx context.Context, sk SubsKey) error {

	sub := sm.subscriptions.Remove(sk.String())
	if sub == nil {
		return gerror.Newf("not found subscription of '%v'", sk)
	}
	if err := sm.streamMgr.DeleteConsumer(ctx, sk.StreamName(), sk.ConsumerName(), sk.TopicName()); err != nil {
		return err
	}
	if err := sub.(*streamConsume).Stop(ctx); err != nil {
		return err
	}
	sm.logger.Infof(ctx, "delete cunsumer <%v> for topic <%v> of stream <%v> ok", sk.ConsumerName(), sk.TopicName(), sk.StreamName())
	return nil
}
