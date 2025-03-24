package nats

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

const (
	consumeMessageDelay = 100 * time.Millisecond
)

// 订阅器管理
type SubscriptionManager struct {
	subscriptions map[string]*subscription
	subMutex      sync.Mutex
	topicBuffers  map[string]chan *pubsub.Message
	bufferMutex   sync.RWMutex
	bufferSize    int // 存放消息的缓存大小
}

type subscription struct {
	cancel context.CancelFunc
}

func newSubscriptionManager(bufferSize int) *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[string]*subscription),
		topicBuffers:  make(map[string]chan *pubsub.Message),
		bufferSize:    bufferSize,
	}
}

// 订阅指定的topic的消息
func (sm *SubscriptionManager) Subscribe(
	ctx context.Context,
	topic string,
	js jetstream.JetStream,
	cfg *Config) (*pubsub.Message, error) {

	metrics.IncrementCounter(ctx, "app_pubsub_subscribe_total_count", "topic", topic)

	if err := sm.validateSubscribePrerequisites(js, cfg); err != nil {
		return nil, err
	}

	sm.subMutex.Lock()

	_, exists := sm.subscriptions[topic]
	if !exists {
		// 获取consumer
		consumer, err := sm.createOrUpdateConsumer(ctx, js, topic, cfg)
		if err != nil {
			sm.subMutex.Unlock()
			return nil, err
		}
		// 注册订阅者
		subCtx, cancel := context.WithCancel(ctx)
		sm.subscriptions[topic] = &subscription{cancel: cancel}
		// 获取主题对应的消息队列缓存
		buffer := sm.getOrCreateBuffer(topic)
		go sm.consumeMessages(subCtx, consumer, topic, buffer, cfg)
	}

	sm.subMutex.Unlock()

	buffer := sm.getOrCreateBuffer(topic)

	select {
	case msg := <-buffer:
		metrics.IncrementCounter(ctx, "app_pubsub_subscribe_success_count", "topic", topic)
		return msg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (*SubscriptionManager) validateSubscribePrerequisites(js jetstream.JetStream, cfg *Config) error {
	if js == nil {
		return errJetStreamNotConfigured
	}

	if cfg.Consumer == "" {
		return errConsumerNotProvided
	}

	return nil
}

func (sm *SubscriptionManager) getOrCreateBuffer(topic string) chan *pubsub.Message {
	sm.bufferMutex.Lock()
	defer sm.bufferMutex.Unlock()

	if buffer, exists := sm.topicBuffers[topic]; exists {
		return buffer
	}

	buffer := make(chan *pubsub.Message, sm.bufferSize)
	sm.topicBuffers[topic] = buffer

	return buffer
}

// 创建或更新一个消费者
func (*SubscriptionManager) createOrUpdateConsumer(
	ctx context.Context, js jetstream.JetStream, topic string, cfg *Config) (jetstream.Consumer, error) {
	// consumerName := fmt.Sprintf("%s_%s", cfg.Consumer, strings.ReplaceAll(topic, ".", "_"))
	consumerName := cfg.Consumer
	cons, err := js.CreateOrUpdateConsumer(ctx, cfg.Stream.Stream, jetstream.ConsumerConfig{
		Durable:       consumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: topic,
		MaxDeliver:    cfg.Stream.MaxDeliver,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckWait:       30 * time.Second,
	})

	return cons, err
}

func (sm *SubscriptionManager) consumeMessages(
	ctx context.Context,
	cons jetstream.Consumer,
	topic string,
	buffer chan *pubsub.Message,
	cfg *Config) {
	// TODO: propagate errors to caller
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := sm.fetchAndProcessMessages(ctx, cons, topic, buffer, cfg); err != nil {
				g.Log().Errorf(ctx, "Error fetching messages for topic %s: %v", topic, err)
			}
		}
	}
}

func (sm *SubscriptionManager) fetchAndProcessMessages(
	ctx context.Context,
	cons jetstream.Consumer,
	topic string,
	buffer chan *pubsub.Message,
	cfg *Config) error {
	msgs, err := cons.Fetch(1, jetstream.FetchMaxWait(cfg.MaxWait))
	if err != nil {
		return sm.handleFetchError(ctx, err, topic)
	}
	return sm.processFetchedMessages(ctx, msgs, topic, buffer)
}

func (*SubscriptionManager) handleFetchError(ctx context.Context, err error, topic string) error {
	if !errors.Is(err, context.DeadlineExceeded) {
		g.Log().Errorf(ctx, "Error fetching messages for topic %s: %v", topic, err)
	}
	time.Sleep(consumeMessageDelay)
	return nil
}

func (sm *SubscriptionManager) processFetchedMessages(
	ctx context.Context,
	msgs jetstream.MessageBatch,
	topic string,
	buffer chan *pubsub.Message) error {
	for msg := range msgs.Messages() {
		pubsubMsg := sm.createPubSubMessage(msg, topic)
		if !sm.sendToBuffer(pubsubMsg, buffer) {
			g.Log().Warningf(ctx, "Message buffer is full for topic %s. Consider increasing buffer size or processing messages faster.", topic)
		}
	}
	return sm.checkBatchError(ctx, msgs, topic)
}

func (*SubscriptionManager) createPubSubMessage(msg jetstream.Msg, topic string) *pubsub.Message {
	pubsubMsg := pubsub.NewMessage(context.Background()) // Pass a context if needed
	pubsubMsg.Topic = topic
	pubsubMsg.Value = msg.Data()
	pubsubMsg.MetaData = msg.Headers()
	pubsubMsg.Committer = &natsCommitter{msg: msg}

	return pubsubMsg
}

// 如果队列满了， 返回false
func (*SubscriptionManager) sendToBuffer(msg *pubsub.Message, buffer chan *pubsub.Message) bool {
	select {
	case buffer <- msg:
		return true
	default:
		return false
	}
}

func (*SubscriptionManager) checkBatchError(ctx context.Context, msgs jetstream.MessageBatch, topic string) error {
	if err := msgs.Error(); err != nil {
		g.Log().Errorf(ctx, "Error in message batch for topic %s: %v", topic, err)
		return err
	}

	return nil
}

// 变比订阅器管理
func (sm *SubscriptionManager) Close() {
	sm.subMutex.Lock()
	for _, sub := range sm.subscriptions {
		sub.cancel()
	}

	sm.subscriptions = make(map[string]*subscription)
	sm.subMutex.Unlock()

	sm.bufferMutex.Lock()
	for _, buffer := range sm.topicBuffers {
		close(buffer)
	}

	sm.topicBuffers = make(map[string]chan *pubsub.Message)

	sm.bufferMutex.Unlock()
}
