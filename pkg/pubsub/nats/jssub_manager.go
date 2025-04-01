package nats

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

func panicRecovery(ctx context.Context, re any) {
	if re == nil {
		return
	}
	g.Log().Error(ctx, re)
}

// 根据消费主题自动生成一个消费者的名字，带有通配符的主题，需要替换通配符
func generateConsumerName(consumer, subject string) string {
	subject = strings.ReplaceAll(subject, ".", "_")
	subject = strings.ReplaceAll(subject, "*", "token")
	subject = strings.ReplaceAll(subject, ">", "tokens")
	return fmt.Sprintf("%s_%s", consumer, subject)
}

const (
	ConsumeMessageDelay = 100 * time.Millisecond
)

type SubType string

var (
	SubTypeJSConsumeNext  SubType = "js-next"
	SubTypeJSConsumeFetch SubType = "js-fetch"
	SubTypeSubAsync       SubType = "sub-async"
	SubTypeSubSync        SubType = "sub-sync"
)

type streamIntf interface {
	createConsumer(ctx context.Context, streamName, consumerName, topicName string) (jetstream.Consumer, error)
	deleteConsumer(ctx context.Context, streamName, consumerName, topicName string) error
}

type closer struct {
	cancel context.CancelFunc
	stop   func()
}
type jsSubscriber struct {
	consumeType                         SubType
	streamName, consumerName, topicName string
	stream                              streamIntf
	close                               closer
	exitNotify                          chan []string
}

// 订阅器管理
type JsSubscriberManager struct {
	subscriptions map[string]*jsSubscriber
	subMutex      sync.Mutex
	exitNotify    chan []string
}

func NewJsSubMgr() *JsSubscriberManager {
	sm := &JsSubscriberManager{
		subscriptions: make(map[string]*jsSubscriber),
		exitNotify:    make(chan []string),
	}
	// 当订阅失败，或stream被删除后，需要删除相关资源
	go func() {
		for key := range sm.exitNotify {
			if len(key) != 0 {
				sm.deleteSubscriber(gctx.New(), key)
			}
		}
	}()
	return sm
}

func (sm *JsSubscriberManager) NewSubscriber(ctx context.Context, stream streamIntf, identity []string, consumeType SubType) error {
	sub := &jsSubscriber{
		stream:       stream,
		streamName:   identity[0],
		consumerName: identity[1],
		topicName:    identity[2],
		consumeType:  consumeType,
		exitNotify:   sm.exitNotify,
	}
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()
	if old, exist := sm.subscriptions[strings.Join(identity, "")]; exist {
		return old.deleteConsumer(ctx)
	}
	sm.subscriptions[strings.Join(identity, "")] = sub
	g.Log().Infof(ctx, "create subscriber ok. streamName: '%v' consumerName: '%v' topicName: '%v'", identity[0], identity[1], identity[2])
	return nil
}

func (sm *JsSubscriberManager) Close(ctx context.Context) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	for _, sub := range sm.subscriptions {
		if err := sub.deleteConsumer(ctx); err != nil {
			return err
		}
	}
	sm.subscriptions = make(map[string]*jsSubscriber)
	return nil
}
func (sm *JsSubscriberManager) deleteSubscriber(ctx context.Context, identity []string) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	if sub, exist := sm.subscriptions[strings.Join(identity, "")]; exist {
		if err := sub.deleteConsumer(ctx); err != nil {
			return err
		}
		sm.subscriptions = nil
	}
	g.Log().Infof(ctx, "delete subscriber ok. streamName:%v consumerName:%v topicName%v", identity[0], identity[1], identity[2])
	return nil
}
func (sm *JsSubscriberManager) DeleteSubscriber(ctx context.Context, identity []string) error {
	return sm.deleteSubscriber(ctx, identity)
}
func (sm *JsSubscriberManager) Subscribe(ctx context.Context, identity []string, handler pubsub.SubscribeFunc) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	if sub, exist := sm.subscriptions[strings.Join(identity, "")]; exist {
		if err := sub.Subscribe(ctx, handler); err != nil {
			return err
		}
	}
	g.Log().Infof(ctx, "delete subscriver ok. streamName:%v consumerName:%v topicName%v", identity[0], identity[1], identity[2])

	return nil
}

// 创建一个消费者
func (s *jsSubscriber) createConsumer(ctx context.Context) (jetstream.Consumer, error) {
	cons, err := s.stream.createConsumer(ctx, s.streamName, s.consumerName, s.topicName)
	if err != nil {
		return nil, gerror.Wrap(err, "createConsumer consumer fail")
	}
	g.Log().Infof(ctx, "create jetstream cunsumer '%v' for topic '%v' of stream '%v' ok", s.consumerName, s.topicName, s.streamName)
	return cons, nil
}

// 删除一个消费者
func (s *jsSubscriber) deleteConsumer(ctx context.Context) error {
	err := s.stream.deleteConsumer(ctx, s.streamName, s.consumerName, s.topicName)
	if err != nil {
		return gerror.Wrap(err, "deleteConsumer consumer fail")
	}
	if s.close.stop != nil {
		s.close.stop()
	}
	if s.close.cancel != nil {
		if s.close.cancel != nil {
			s.close.cancel()
		}
	}
	g.Log().Infof(ctx, "delete cunsumer <%v> for topic <%v> of stream <%v> ok", s.consumerName, s.topicName, s.streamName)
	return nil
}
func (s *jsSubscriber) Subscribe(ctx context.Context, handler pubsub.SubscribeFunc) error {
	switch s.consumeType {
	case SubTypeJSConsumeNext:
		return s.subscribeByNext(ctx, handler)
	}
	return nil
}

// 订阅指定的topic的消息
func (s *jsSubscriber) subscribeByNext(
	ctx context.Context,
	handler pubsub.SubscribeFunc,
) error {
	metrics.IncrementCounter(ctx, metrics.NatsSubscribeTotalCount, "topic", s.topicName)

	// 获取consumer
	consumer, err := s.createConsumer(ctx)
	if err != nil {
		return err
	}
	// 获取消息迭代器
	iter, err := consumer.Messages(jetstream.PullMaxMessages(1))
	if err != nil {
		return gerror.Wrap(err, "get consumer msg iter fail")
	}
	// 注册订阅者
	subCtx, cancel := context.WithCancel(ctx)
	s.close = closer{
		cancel: cancel,
		stop:   iter.Stop,
	}
	defer func() {
		s.exitNotify <- []string{s.streamName, s.consumerName, s.topicName}
	}()
	g.Log().Infof(ctx, "consumer '%v' ready to consume msg for '%v' of '%v'", s.consumerName, s.topicName, s.streamName)
	// 获取主题对应的消息队列缓存
	for {
		select {
		case <-subCtx.Done():
			g.Log().Infof(subCtx, "stream '%v' consumer '%v' closer cancelled fot topic '%v'", s.streamName, s.consumerName, s.topicName)
			return nil
		default:
			msg, err := iter.Next()
			g.Log().Warningf(ctx, "js consume:%s", msg.Data())
			if err != nil {
				if !errors.Is(err, jetstream.ErrNoHeartbeat) {
					g.Log().Warningf(ctx, "consumer '%v' fetching messages for topic '%s' of '%v': %v", s.consumerName, s.topicName, s.streamName, err)
				} else {
					g.Log().Errorf(ctx, "consumer '%v' fetching messages for topic '%s' of '%v' fail: %v", s.consumerName, s.topicName, s.streamName, err)
					return nil
				}
				time.Sleep(ConsumeMessageDelay)
				continue
			}
			pubsubmsg := s.createPubSubMessage(msg, s.topicName)
			err = func() error {
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
				g.Log().Errorf(ctx, "consumer '%v' error in handler for topic '%s': %v", s.consumerName, s.topicName, err)
				continue
			}
			// 处理完成
			if pubsubmsg.Committer != nil {
				pubsubmsg.Commit(ctx)
			}
			continue
		}
	}
}

func (sm *jsSubscriber) createPubSubMessage(msg jetstream.Msg, topic string) *pubsub.Message {
	pubsubMsg := pubsub.NewMessage() // Pass a context if needed
	pubsubMsg.Topic = topic
	pubsubMsg.Value = msg.Data()
	pubsubMsg.MetaData = msg.Headers()
	pubsubMsg.Committer = &jsCommitter{msg: msg}
	return pubsubMsg
}
