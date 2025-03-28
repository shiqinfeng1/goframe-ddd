package nats

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

const (
	ConsumeMessageDelay = 100 * time.Millisecond
)

// 订阅器管理
type SubscriptionManager struct {
	subscriptions map[string]*subscription
	subMutex      sync.Mutex
}

type subscription struct {
	cancel      context.CancelFunc
	msgIterStop func()
}

func newSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[string]*subscription),
	}
}

func (*SubscriptionManager) validateSubscribePrerequisites(ctx context.Context, js jetstream.JetStream, cfg *Config) error {
	if js == nil {
		return errJetStreamNotConfigured
	}
	if cfg.ConsumerName == "" {
		return errConsumerNotProvided
	}
	_, err := js.Stream(ctx, cfg.Stream.Name)
	if err != nil {
		return errGetStream
	}
	return nil
}

// 创建或更新一个消费者
func (*SubscriptionManager) createOrUpdateConsumer(
	ctx context.Context, js jetstream.JetStream, topic string, cfg *Config,
) (jetstream.Consumer, error) {
	consumerName := generateConsumerName(cfg.ConsumerName, topic)
	cons, err := js.CreateOrUpdateConsumer(ctx, cfg.Stream.Name, jetstream.ConsumerConfig{
		Durable:       consumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: topic,
		MaxDeliver:    cfg.Stream.MaxDeliver,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckWait:       30 * time.Second,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "create consumer fail")
	}
	g.Log().Infof(ctx, "create or updateing cunsumer <%v> for stream <%v> ok", consumerName, cfg.Stream.Name)
	return cons, nil
}

// 订阅指定的topic的消息
func (sm *SubscriptionManager) Subscribe(
	ctx context.Context,
	topic string,
	js jetstream.JetStream,
	cfg *Config,
	handler pubsub.SubscribeFunc,
) error {
	metrics.IncrementCounter(ctx, metrics.NatsSubscribeTotalCount, "topic", topic)

	if err := sm.validateSubscribePrerequisites(ctx, js, cfg); err != nil {
		return err
	}

	sm.subMutex.Lock()

	_, exists := sm.subscriptions[topic]
	if exists {
		sm.subMutex.Unlock()
		return gerror.Newf("consumer %v is already subscribed topic %v", cfg.ConsumerName, topic)
	}
	// 获取consumer
	consumer, err := sm.createOrUpdateConsumer(ctx, js, topic, cfg)
	if err != nil {
		sm.subMutex.Unlock()
		return err
	}
	// 获取消息迭代器
	iter, err := consumer.Messages(jetstream.PullMaxMessages(1))
	if err != nil {
		sm.subMutex.Unlock()
		return gerror.Wrap(err, "get consumer msg iter fail")
	}
	// 注册订阅者
	subCtx, cancel := context.WithCancel(ctx)
	sm.subscriptions[topic] = &subscription{
		cancel:      cancel,
		msgIterStop: iter.Stop,
	}
	sm.subMutex.Unlock()
	// 获取主题对应的消息队列缓存
	sm.consumeMessages(subCtx, iter, topic, cfg, handler)
	return nil
}

func panicRecovery(ctx context.Context, re any) {
	if re == nil {
		return
	}
	g.Log().Error(ctx, re)
}

func (sm *SubscriptionManager) consumeMessages(
	ctx context.Context,
	msgIter jetstream.MessagesContext,
	topic string,
	cfg *Config,
	handler pubsub.SubscribeFunc,
) {
	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "consumer %v subscription cancelled fot topic %v", cfg.ConsumerName, topic)
			return
		default:
			g.Log().Infof(ctx, "%v ready to consume msg for %v", cfg.ConsumerName, topic)
			msg, err := msgIter.Next()
			if err != nil {
				if !errors.Is(err, jetstream.ErrNoMessages) {
					g.Log().Warningf(ctx, "consumer %v error fetching messages for topic %s: %v", cfg.ConsumerName, topic, err)
				} else {
					g.Log().Errorf(ctx, "consumer %v fetching messages for topic %s fail: %v", cfg.ConsumerName, topic, err)
				}
				time.Sleep(ConsumeMessageDelay)
				continue
			}
			pubsubmsg := sm.createPubSubMessage(msg, topic)
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
				g.Log().Errorf(ctx, "consumer %v error in handler for topic '%s': %v", cfg.ConsumerName, topic, err)
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
func (sm *SubscriptionManager) createPubSubMessage(msg jetstream.Msg, topic string) *pubsub.Message {
	pubsubMsg := pubsub.NewMessage(context.Background()) // Pass a context if needed
	pubsubMsg.Topic = topic
	pubsubMsg.Value = msg.Data()
	pubsubMsg.MetaData = msg.Headers()
	pubsubMsg.Committer = &natsCommitter{msg: msg}
	return pubsubMsg
}

// 变比订阅器管理
func (sm *SubscriptionManager) Close() {
	sm.subMutex.Lock()
	defer sm.subMutex.Unlock()

	for _, sub := range sm.subscriptions {
		sub.msgIterStop()
		sub.cancel()
	}

	sm.subscriptions = make(map[string]*subscription)
}
