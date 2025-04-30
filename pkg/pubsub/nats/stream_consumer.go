package nats

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

// 订阅器管理
type StreamConsumer struct {
	logger        pubsub.Logger
	subscriptions map[string]*streamConsume
	streamMgr     *JetStreamWrapper
	subMutex      sync.Mutex
	exitNotify    chan SubsKey // 同步模式下。订阅退出后，各个consume通知本consumer删除相关记录
}

func NewStreamConsumer(logger pubsub.Logger) *StreamConsumer {
	sm := &StreamConsumer{
		logger:        logger,
		subscriptions: make(map[string]*streamConsume),
		subMutex:      sync.Mutex{},
		exitNotify:    make(chan SubsKey),
	}
	// 当订阅失败，或stream被删除后，需要删除相关资源
	go func() {
		for key := range sm.exitNotify {
			sm.DeleteConsumer(gctx.New(), key)
		}
	}()
	return sm
}

func (sm *StreamConsumer) AddConsume(
	ctx context.Context,
	st SubType,
	sk SubsKey,
	c jetstream.Consumer,
	handler ConsumeFunc,
	en chan SubsKey) error {

	ssub := NewStreamConsume(sm.logger, st, sk, c, handler, en)

	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()
	if _, exist := sm.subscriptions[sk.String()]; exist {
		return gerror.Newf("stream '%v' topic '%v' is already be consumed by '%v'", sk.StreamName(), sk.TopicName(), sk.ConsumerName())
	}
	sm.subscriptions[sk.String()] = ssub
	if err := ssub.start(ctx); err != nil {
		return err
	}
	sm.logger.Infof(ctx, "stream '%v' create consumer '%v' of topic '%v' ok", sk.StreamName(), sk.ConsumerName(), sk.TopicName())
	return nil
}

func (sm *StreamConsumer) Close(ctx context.Context) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()
	for _, sub := range sm.subscriptions {
		if err := sm.streamMgr.DeleteConsumer(ctx, sub.subsKey.StreamName(), sub.subsKey.ConsumerName(), sub.subsKey.TopicName()); err != nil {
			return err
		}
		if err := sub.Stop(ctx); err != nil {
			return err
		}
	}
	sm.subscriptions = make(map[string]*streamConsume)
	close(sm.exitNotify)
	return nil
}

func (sm *StreamConsumer) DeleteConsumer(ctx context.Context, sk SubsKey) error {

	sm.subMutex.Lock()
	sub, exist := sm.subscriptions[sk.String()]
	if !exist {
		sm.subMutex.Unlock()
		return gerror.Newf("not found subscription of '%v'", sk)
	}
	sm.subscriptions = nil
	sm.subMutex.Unlock()

	if err := sm.streamMgr.DeleteConsumer(ctx, sk.StreamName(), sk.ConsumerName(), sk.TopicName()); err != nil {
		return err
	}
	if err := sub.Stop(ctx); err != nil {
		return err
	}
	sm.logger.Infof(ctx, "delete cunsumer <%v> for topic <%v> of stream <%v> ok", sk.ConsumerName(), sk.TopicName(), sk.StreamName())
	return nil
}
