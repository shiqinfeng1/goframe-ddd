package nats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewSubscriptionManager(t *testing.T) {
	sm := newSubscriptionManager()
	assert.NotNil(t, sm)
	assert.NotNil(t, sm.subscriptions)
}

func TestSubscriptionManager_validateSubscribePrerequisites(t *testing.T) {
	sm := newSubscriptionManager()
	mockJS := NewMockJetStream(gomock.NewController(t))
	mockJS.EXPECT().Stream(gomock.Any(), gomock.Any()).Return(nil, nil)

	cfg := &Config{ConsumerName: "test-consumer"}

	err := sm.validateSubscribePrerequisites(t.Context(), mockJS, cfg)
	require.NoError(t, err)

	err = sm.validateSubscribePrerequisites(t.Context(), nil, cfg)
	assert.Equal(t, errJetStreamNotConfigured, err)

	err = sm.validateSubscribePrerequisites(t.Context(), mockJS, &Config{})
	assert.Equal(t, errConsumerNotProvided, err)
}

func TestSubscriptionManager_createOrUpdateConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)
	mockConsumer := NewMockConsumer(ctrl)

	sm := newSubscriptionManager()
	cfg := &Config{
		ConsumerName: "test-consumer",
		Stream: StreamConfig{
			Name:       "test-stream",
			MaxDeliver: 3,
		},
	}

	ctx := context.Background()
	topic := "test.topic"

	mockJS.EXPECT().CreateOrUpdateConsumer(ctx, cfg.Stream.Name, gomock.Any()).Return(mockConsumer, nil)

	consumer, err := sm.createOrUpdateConsumer(ctx, mockJS, topic, cfg)
	require.NoError(t, err)
	assert.Equal(t, mockConsumer, consumer)
}

func TestSubscriptionManager_Close(t *testing.T) {
	sm := newSubscriptionManager()
	topic := "test.topic"

	// Create a subscription and buffer
	ctx, cancel := context.WithCancel(context.Background())
	sm.subscriptions[topic] = &subscription{
		cancel:      cancel,
		msgIterStop: func() {},
	}

	sm.Close()

	assert.Empty(t, sm.subscriptions)

	// Check that the context was canceled
	if ctx.Err() == nil {
		t.Fatal("Context was not canceled")
	}
}
