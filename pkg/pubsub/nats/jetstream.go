package natsclient

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go/jetstream"
	pubsub "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

// JetStreamWrapper is a manager for jStream streams.
type JetStreamWrapper struct {
	logger pubsub.Logger
	js     jetstream.JetStream
}

// newStreamManager creates a new JetStreamWrapper.
func NewJetStreamWrapper(logger pubsub.Logger, js jetstream.JetStream) *JetStreamWrapper {
	return &JetStreamWrapper{
		logger: logger,
		js:     js,
	}
}

func (sm *JetStreamWrapper) validateStream(_ context.Context, name string, subjects []string) error {
	if !sm.js.Conn().IsConnected() {
		return gerror.New("nats not connected")
	}
	if len(subjects) == 0 {
		return gerror.New("subjects is nil")
	}
	if name == "" {
		return gerror.New("stream name is nil")
	}
	return nil
}

// CreateStream creates a new jStream stream.
func (sm *JetStreamWrapper) CreateOrUpdateStream(ctx context.Context, name string, subjects []string) error {
	if err := sm.validateStream(ctx, name, subjects); err != nil {
		return err
	}
	jsCfg := jetstream.StreamConfig{
		Name:       name,
		Subjects:   subjects,
		Storage:    jetstream.FileStorage, // 默认文件存储
		MaxMsgSize: 10 * 1024 * 1024,
		Retention:  jetstream.InterestPolicy, // 如果有多个消费者订阅了相同的主题，等每个消费者都消费确认后删除消息
	}

	_, err := sm.js.CreateOrUpdateStream(ctx, jsCfg)
	if err != nil {
		return gerror.Wrapf(err, "create or update stream fail")
	}
	sm.logger.Debugf(ctx, "create or update stream ok. stream-name=%v, subjects=%+v", name, subjects)
	return nil
}

// CreateStream creates a new jStream stream.
func (sm *JetStreamWrapper) CreateStream(ctx context.Context, name string, subjects []string) error {
	if err := sm.validateStream(ctx, name, subjects); err != nil {
		return err
	}
	// todo：根据需求需要更详细配置
	jsCfg := jetstream.StreamConfig{
		Name:       name,
		Subjects:   subjects,
		Storage:    jetstream.FileStorage,    // 默认文件存储
		MaxMsgSize: 10 * 1024 * 1024,         // 10M
		Retention:  jetstream.InterestPolicy, // 如果有多个消费者订阅了相同的主题，每个消费者都可能接收到相同的消息
	}
	_, err := sm.js.CreateStream(ctx, jsCfg)
	if err != nil {
		return gerror.Wrapf(err, "failed to create stream")
	}
	sm.logger.Debugf(ctx, "creating stream '%s' ok of subjects: '%+v'", name, subjects)

	return nil
}

// DeleteStream deletes a jStream stream.
func (sm *JetStreamWrapper) DeleteStream(ctx context.Context, name string) error {
	sm.logger.Debugf(ctx, "deleteing stream '%s'", name)

	err := sm.js.DeleteStream(ctx, name)
	if err != nil {
		if errors.Is(err, jetstream.ErrStreamNotFound) {
			sm.logger.Debugf(ctx, "stream '%s' not found, considering delete successful", name)
			return nil // If the stream doesn't exist, we consider it a success
		}
		return gerror.Wrapf(err, "failed to delete stream '%s'", name)
	}
	sm.logger.Debugf(ctx, "successfully deleted stream '%s'", name)
	return nil
}

// GetStream gets a jStream stream.
func (sm *JetStreamWrapper) GetStream(ctx context.Context, name string) (jetstream.Stream, error) {
	stream, err := sm.js.Stream(ctx, name)
	if err != nil {
		if errors.Is(err, jetstream.ErrStreamNotFound) {
			return nil, gerror.Wrapf(err, "stream %s not found", name)
		}
		return nil, gerror.Wrapf(err, "failed to get stream %s", name)
	}
	info, err := stream.Info(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "get stream info fail")
	}
	sm.logger.Debugf(ctx, "getting stream info ok: %+v", info)

	return stream, nil
}

// 根据消费主题自动生成一个消费者的名字，带有通配符的主题，需要替换通配符
func genConsumerName(consumer, subject string) string {
	subject = strings.ReplaceAll(subject, ".", "_")
	subject = strings.ReplaceAll(subject, "*", "token")
	subject = strings.ReplaceAll(subject, ">", "tokens")
	return fmt.Sprintf("%s_%s", consumer, subject)
}

func (sm *JetStreamWrapper) CreateOrUpdateConsumer(ctx context.Context, streamName, consumerName, subject string) (jetstream.Consumer, error) {
	cons, err := sm.js.CreateOrUpdateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Name:          genConsumerName(consumerName, subject),
		Durable:       genConsumerName(consumerName, subject),
		AckPolicy:     jetstream.AckExplicitPolicy, //AckExplicitPolicy,
		FilterSubject: subject,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckWait:       30 * time.Second, // 业务处理消息的最长时间，如果该时间内没有回复ack，将重推送该消息
		// MaxAckPending: 1000,                // 最多为回复ack的消息数量，如果到达上限，服务端将停止推送
	})
	if err != nil {
		return nil, gerror.Wrapf(err, "create or update consumer fail. stream-name=%v, consumer-name=%v, subject=%v", streamName, consumerName, subject)
	}
	sm.logger.Debugf(ctx, "creating or updating consumer ok. stream-name=%v, consumer-name=%v, subject=%v", streamName, consumerName, subject)
	return cons, nil
}

func (sm *JetStreamWrapper) DeleteConsumer(ctx context.Context, streamName, consumerName, subject string) error {
	if !sm.js.Conn().IsConnected() {
		return gerror.Newf("delete consumer fail: nats not connected")
	}
	err := sm.js.DeleteConsumer(ctx, streamName, genConsumerName(consumerName, subject))
	if err != nil {
		return gerror.Wrapf(err, "delete consumer fail. stream-name=%v, consumer-name=%v, subject=%v", streamName, consumerName, subject)
	}
	sm.logger.Debugf(ctx, "deleting consumer ok. stream-name=%v, consumer-name=%v, subject=%v", streamName, consumerName, subject)
	return nil
}
