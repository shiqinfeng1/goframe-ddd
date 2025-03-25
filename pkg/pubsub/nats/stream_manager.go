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
func newStreamManager(js jetstream.JetStream) *StreamManager {
	return &StreamManager{
		js: js,
	}
}

// CreateStream creates a new jStream stream.
func (sm *StreamManager) CreateStream(ctx context.Context, cfg StreamConfig) error {
	g.Log().Debugf(ctx, "creating stream %s", cfg.Stream)
	// todo：根据需求需要更详细配置
	jsCfg := jetstream.StreamConfig{
		Name:     cfg.Stream,
		Subjects: cfg.Subjects,
		MaxBytes: cfg.MaxBytes,
	}

	_, err := sm.js.CreateStream(ctx, jsCfg)
	if err != nil {
		return gerror.Wrap(err, "failed to create stream")
	}

	return nil
}

// DeleteStream deletes a jStream stream.
func (sm *StreamManager) DeleteStream(ctx context.Context, name string) error {
	g.Log().Debugf(ctx, "deleting stream %s", name)

	err := sm.js.DeleteStream(ctx, name)
	if err != nil {
		if errors.Is(err, jetstream.ErrStreamNotFound) {
			g.Log().Debugf(ctx, "stream %s not found, considering delete successful", name)
			return nil // If the stream doesn't exist, we consider it a success
		}
		return gerror.Wrapf(err, "failed to delete stream %s", name)
	}
	g.Log().Debugf(ctx, "successfully deleted stream %s", name)
	return nil
}

// CreateOrUpdateStream creates or updates a jStream stream.
func (sm *StreamManager) CreateOrUpdateStream(ctx context.Context, cfg *jetstream.StreamConfig) (jetstream.Stream, error) {
	g.Log().Debugf(ctx, "creating or updating stream %s", cfg.Name)

	stream, err := sm.js.CreateOrUpdateStream(ctx, *cfg)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to create or update stream")
	}
	return stream, nil
}

// GetStream gets a jStream stream.
func (sm *StreamManager) GetStream(ctx context.Context, name string) (jetstream.Stream, error) {
	g.Log().Debugf(ctx, "getting stream %s", name)

	stream, err := sm.js.Stream(ctx, name)
	if err != nil {
		if errors.Is(err, jetstream.ErrStreamNotFound) {
			return nil, gerror.Wrapf(err, "stream %s not found", name)
		}
		return nil, gerror.Wrapf(err, "failed to get stream %s", name)
	}

	return stream, nil
}
