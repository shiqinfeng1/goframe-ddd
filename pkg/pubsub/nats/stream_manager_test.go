package nats

import (
	"context"
	"errors"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewStreamManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	sm := newStreamManager(mockJS)

	assert.NotNil(t, sm)
	assert.Equal(t, mockJS, sm.js)
}

func TestStreamManager_CreateStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	sm := newStreamManager(mockJS)

	ctx := context.Background()
	cfg := StreamConfig{
		Name:     "test-stream",
		Subjects: []string{"test.subject"},
	}

	mockJS.EXPECT().CreateOrUpdateStream(ctx, gomock.Any()).Return(nil, nil)

	err := sm.CreateStream(ctx, cfg)
	require.NoError(t, err)
}

func TestStreamManager_CreateStream_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	sm := newStreamManager(mockJS)

	ctx := context.Background()
	cfg := StreamConfig{
		Name:     "test-stream",
		Subjects: []string{"test.subject"},
	}

	expectedErr := errCreateStream
	mockJS.EXPECT().CreateOrUpdateStream(ctx, gomock.Any()).Return(nil, expectedErr)

	err := sm.CreateStream(ctx, cfg)
	require.Error(t, err)
	assert.Equal(t, expectedErr, errors.Unwrap(err))
}

func TestStreamManager_DeleteStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	sm := newStreamManager(mockJS)

	ctx := context.Background()
	streamName := "test-stream"

	mockJS.EXPECT().DeleteStream(ctx, streamName).Return(nil)

	err := sm.DeleteStream(ctx, streamName)
	require.NoError(t, err)
}

func TestStreamManager_DeleteStream_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	sm := newStreamManager(mockJS)

	ctx := context.Background()
	streamName := "test-stream"

	mockJS.EXPECT().DeleteStream(ctx, streamName).Return(jetstream.ErrStreamNotFound)

	err := sm.DeleteStream(ctx, streamName)
	require.NoError(t, err)
}

func TestStreamManager_DeleteStream_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	sm := newStreamManager(mockJS)

	ctx := context.Background()
	streamName := "test-stream"

	expectedErr := errDeleteStream
	mockJS.EXPECT().DeleteStream(ctx, streamName).Return(expectedErr)

	err := sm.DeleteStream(ctx, streamName)
	require.Error(t, err)
	assert.Equal(t, expectedErr, errors.Unwrap(err))
}

func TestStreamManager_GetStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	sm := newStreamManager(mockJS)

	ctx := context.Background()
	streamName := "test-stream"

	mockStream := NewMockStream(ctrl)
	mockJS.EXPECT().Stream(ctx, streamName).Return(mockStream, nil)

	stream, err := sm.GetStream(ctx, streamName)
	require.NoError(t, err)
	assert.Equal(t, mockStream, stream)
}

func TestStreamManager_GetStream_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	sm := newStreamManager(mockJS)

	ctx := context.Background()
	streamName := "test-stream"

	mockJS.EXPECT().Stream(ctx, streamName).Return(nil, jetstream.ErrStreamNotFound)

	stream, err := sm.GetStream(ctx, streamName)
	require.Error(t, err)
	assert.Nil(t, stream)
	assert.Equal(t, jetstream.ErrStreamNotFound, errors.Unwrap(err))
}

func TestStreamManager_GetStream_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	sm := newStreamManager(mockJS)

	ctx := context.Background()
	streamName := "test-stream"

	expectedErr := errGetStream
	mockJS.EXPECT().Stream(ctx, streamName).Return(nil, expectedErr)

	stream, err := sm.GetStream(ctx, streamName)
	require.Error(t, err)
	assert.Nil(t, stream)
	assert.Equal(t, expectedErr, errors.Unwrap(err))
}
