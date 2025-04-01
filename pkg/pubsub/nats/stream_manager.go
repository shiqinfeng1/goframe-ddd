package nats

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go/jetstream"
)

// StreamManager is a manager for jStream streams.
type StreamManager struct {
	js jetstream.JetStream
}

// newStreamManager creates a new StreamManager.
func NewStreamManager(js jetstream.JetStream) *StreamManager {
	return &StreamManager{
		js: js,
	}
}

// CreateStream creates a new jStream stream.
func (sm *StreamManager) CreateOrUpdateStream(ctx context.Context, name string, subjects []string) error {
	if len(subjects) == 0 {
		return errSubjectsNotProvided
	}
	if name == "" {
		return errStreamNotProvided
	}
	// todo：根据需求需要更详细配置
	jsCfg := jetstream.StreamConfig{
		Name:      name,
		Subjects:  subjects,
		Storage:   jetstream.FileStorage,    // 默认文件存储
		Retention: jetstream.InterestPolicy, // 如果有多个消费者订阅了相同的主题，每个消费者都可能接收到相同的消息
	}

	_, err := sm.js.CreateOrUpdateStream(ctx, jsCfg)
	if err != nil {
		return gerror.Wrapf(err, "failed to create stream")
	}
	g.Log().Debugf(ctx, "creating or updating stream %s ok of subkects:%+v", name, subjects)

	return nil
}

// CreateStream creates a new jStream stream.
func (sm *StreamManager) CreateStream(ctx context.Context, name string, subjects []string) error {
	if len(subjects) == 0 {
		return errSubjectsNotProvided
	}
	if name == "" {
		return errStreamNotProvided
	}
	// todo：根据需求需要更详细配置
	jsCfg := jetstream.StreamConfig{
		Name:      name,
		Subjects:  subjects,
		Storage:   jetstream.FileStorage,    // 默认文件存储
		Retention: jetstream.InterestPolicy, // 如果有多个消费者订阅了相同的主题，每个消费者都可能接收到相同的消息
	}

	_, err := sm.js.CreateStream(ctx, jsCfg)
	if err != nil {
		return gerror.Wrapf(err, "failed to create stream")
	}
	g.Log().Debugf(ctx, "creating stream '%s' ok of subjects: '%+v'", name, subjects)

	return nil
}

// DeleteStream deletes a jStream stream.
func (sm *StreamManager) DeleteStream(ctx context.Context, name string) error {
	g.Log().Debugf(ctx, "deleting stream '%s'", name)

	err := sm.js.DeleteStream(ctx, name)
	if err != nil {
		if errors.Is(err, jetstream.ErrStreamNotFound) {
			g.Log().Debugf(ctx, "stream '%s' not found, considering delete successful", name)
			return nil // If the stream doesn't exist, we consider it a success
		}
		return gerror.Wrapf(err, "failed to delete stream '%s'", name)
	}
	g.Log().Debugf(ctx, "successfully deleted stream '%s'", name)
	return nil
}

// GetStream gets a jStream stream.
func (sm *StreamManager) GetStream(ctx context.Context, name string) (jetstream.Stream, error) {

	stream, err := sm.js.Stream(ctx, name)
	if err != nil {
		if errors.Is(err, jetstream.ErrStreamNotFound) {
			return nil, gerror.Wrapf(err, "stream %s not found", name)
		}
		return nil, gerror.Wrapf(err, "failed to get stream %s", name)
	}
	info, _ := stream.Info(ctx)
	g.Log().Debugf(ctx, "getting stream info ok: %v", info)

	return stream, nil
}

// GetJetStreamStatus returns the status of the jStream connection.
func (sm *StreamManager) GetJetStreamStatus(ctx context.Context) (string, error) {
	_, err := sm.js.AccountInfo(ctx)
	if err != nil {
		return jetStreamStatusError, err
	}

	return jetStreamStatusOK, nil
}

func (sm *StreamManager) createConsumer(ctx context.Context, streamName, consumerName, subject string) (jetstream.Consumer, error) {
	cons, err := sm.js.CreateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Durable:       generateConsumerName(consumerName, subject),
		AckPolicy:     jetstream.AckExplicitPolicy, //AckExplicitPolicy,
		FilterSubject: subject,
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to create consumer for stream %v", streamName)
	}

	return cons, nil

}
func (sm *StreamManager) deleteConsumer(ctx context.Context, streamName, consumerName, subject string) error {
	err := sm.js.DeleteConsumer(ctx, streamName, generateConsumerName(consumerName, subject))
	if err != nil {
		return gerror.Wrapf(err, "failed to delete consumer for stream %v", streamName)
	}
	return nil
}
