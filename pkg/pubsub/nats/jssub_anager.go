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
	deleteConsumer(ctx context.Context, streamName, consumerName string) error
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
	exitNotify                          chan string
}

// 订阅器管理
type JsSubscriberManager struct {
	subscriptions map[string]*jsSubscriber
	subMutex      sync.Mutex
	exitNotify    chan string
}

func NewJsSubscriberManager() *JsSubscriberManager {
	sm := &JsSubscriberManager{
		subscriptions: make(map[string]*jsSubscriber),
		exitNotify:    make(chan string),
	}
	// 当订阅失败，或stream被删除后，需要删除相关资源
	go func() {
		for {
			select {
			case key, ok := <-sm.exitNotify:
				if ok {
					sm.deleteSubscriber(gctx.New(), key)
				}
			}
		}
	}()
	return sm
}

func (sm *JsSubscriberManager) NewSubscriber(stream streamIntf, streamName, consumerName, topicName string, consumeType SubType) *jsSubscriber {
	sub := &jsSubscriber{
		stream:      stream,
		consumeType: consumeType,
		exitNotify:  sm.exitNotify,
	}
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()
	sm.subscriptions[streamName+consumerName+topicName] = sub
	return sub
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
func (sm *JsSubscriberManager) deleteSubscriber(ctx context.Context, key string) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	if sub, exist := sm.subscriptions[key]; exist {
		if err := sub.deleteConsumer(ctx); err != nil {
			return err
		}
		sm.subscriptions = nil
	}
	return nil
}
func (sm *JsSubscriberManager) DeleteSubscriber(ctx context.Context, streamName, consumerName, topicName string) error {
	return sm.deleteSubscriber(ctx, streamName+consumerName+topicName)
}
func (sm *JsSubscriberManager) Subscribe(ctx context.Context, streamName, consumerName, topicName string, handler pubsub.SubscribeFunc) error {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	if sub, exist := sm.subscriptions[streamName+consumerName+topicName]; exist {
		if err := sub.Subscribe(ctx, handler); err != nil {
			return err
		}
	}
	return nil
}

// 创建一个消费者
func (s *jsSubscriber) createConsumer(ctx context.Context, streamName, consumerName, topicName string) (jetstream.Consumer, error) {
	consumerName2 := generateConsumerName(consumerName, topicName)
	cons, err := s.stream.createConsumer(ctx, streamName, consumerName2, topicName)
	if err != nil {
		return nil, gerror.Wrap(err, "createConsumer consumer fail")
	}
	g.Log().Infof(ctx, "createConsumer cunsumer <%v> for topic <%v> ok", consumerName2, topicName)
	return cons, nil
}

// 删除一个消费者
func (s *jsSubscriber) deleteConsumer(ctx context.Context) error {
	consumerName2 := generateConsumerName(s.consumerName, s.topicName)
	err := s.stream.deleteConsumer(ctx, s.streamName, consumerName2)
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
	g.Log().Infof(ctx, "deleteConsumer cunsumer <%v> for topic <%v> ok", s.consumerName, s.topicName)
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
	consumer, err := s.createConsumer(ctx, s.streamName, s.consumerName, s.topicName)
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
		s.exitNotify <- s.streamName + s.consumerName + s.topicName
	}()
	// 获取主题对应的消息队列缓存
	for {
		select {
		case <-subCtx.Done():
			g.Log().Infof(subCtx, "stream %v consumer %v closer cancelled fot topic %v", s.streamName, s.consumerName, s.topicName)
			return nil
		default:
			g.Log().Infof(ctx, "%v ready to consume msg for %v of %v", s.consumerName, s.topicName, s.streamName)
			msg, err := iter.Next()
			if err != nil {
				if !errors.Is(err, jetstream.ErrNoHeartbeat) {
					g.Log().Warningf(ctx, "consumer %v fetching messages for topic %s of %v: %v", s.consumerName, s.topicName, s.streamName, err)
				} else {
					g.Log().Errorf(ctx, "consumer %v fetching messages for topic %s of %v fail: %v", s.consumerName, s.topicName, s.streamName, err)
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
				g.Log().Errorf(ctx, "consumer %v error in handler for topic '%s': %v", s.consumerName, s.topicName, err)
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
