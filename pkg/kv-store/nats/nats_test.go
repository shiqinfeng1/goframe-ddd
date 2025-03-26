package nats

import (
	"context"
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
	jetstream "github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	errFailedToSet      = errors.New("failed to set")
	errConnectionFailed = errors.New("connection failed")
)

func Test_ClientSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKV := NewMockKeyValue(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}
	mockKV.EXPECT().
		Put(gomock.Any(), "test_key", []byte("test_value")).
		Return(uint64(1), nil)

	cl := Client{
		kv:      mockKV,
		configs: configs,
	}

	err := cl.Set(context.Background(), "test_key", "test_value")
	require.NoError(t, err)
}

func Test_ClientSetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKV := NewMockKeyValue(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}

	mockKV.EXPECT().
		Put(gomock.Any(), "test_key", []byte("test_value")).
		Return(uint64(0), errFailedToSet)

	cl := Client{
		kv:      mockKV,
		configs: configs,
	}

	err := cl.Set(context.Background(), "test_key", "test_value")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set key-value pair")
}

func Test_ClientGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKV := NewMockKeyValue(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}
	//
	mockEntry := NewMockKeyValueEntry(ctrl)
	mockEntry.EXPECT().Value().Return([]byte("test_value"))
	mockKV.EXPECT().
		Get(gomock.Any(), "test_key").
		Return(mockEntry, nil)

	cl := Client{
		kv:      mockKV,
		configs: configs,
	}

	val, err := cl.Get(context.Background(), "test_key")
	require.NoError(t, err)
	assert.Equal(t, "test_value", val)
}

func Test_ClientGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKV := NewMockKeyValue(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}

	mockKV.EXPECT().
		Get(gomock.Any(), "nonexistent_key").
		Return(nil, nats.ErrKeyNotFound)

	cl := Client{
		kv:      mockKV,
		configs: configs,
	}

	val, err := cl.Get(t.Context(), "nonexistent_key")
	require.Error(t, err)
	assert.Empty(t, val)
	assert.Contains(t, err.Error(), "key not found")
}

func Test_ClientDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKV := NewMockKeyValue(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}

	mockKV.EXPECT().
		Delete(gomock.Any(), "test_key").
		Return(nil)

	cl := Client{
		kv:      mockKV,
		configs: configs,
	}

	err := cl.Delete(context.Background(), "test_key")
	require.NoError(t, err)
}

func Test_ClientDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKV := NewMockKeyValue(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}

	mockKV.EXPECT().
		Delete(gomock.Any(), "nonexistent_key").
		Return(nats.ErrKeyNotFound)

	cl := Client{
		kv:      mockKV,
		configs: configs,
	}

	err := cl.Delete(context.Background(), "nonexistent_key")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "key not found")
}

func Test_ClientHealthCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}

	mockJS.EXPECT().
		AccountInfo(gomock.Any()).
		Return(&jetstream.AccountInfo{}, nil)

	cl := Client{
		js:      mockJS,
		configs: configs,
	}

	health := cl.Health(context.Background())
	assert.Equal(t, "UP", health.Status)
	assert.Equal(t, configs.Server, health.Details["url"])
	assert.Equal(t, configs.Bucket, health.Details["bucket"])
}

func Test_ClientHealthCheckFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJetStream(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}

	mockJS.EXPECT().
		AccountInfo(gomock.Any()).
		Return(nil, errConnectionFailed)

	cl := Client{
		js:      mockJS,
		configs: configs,
	}

	health := cl.Health(context.Background())
	assert.Equal(t, "DOWN", health.Status)
	assert.Equal(t, configs.Server, health.Details["url"])
	assert.Equal(t, configs.Bucket, health.Details["bucket"])
}
