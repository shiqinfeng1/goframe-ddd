package nats

import (
	"context"
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
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
		Put("test_key", []byte("test_value")).
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
		Put("test_key", []byte("test_value")).
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

	mockEntry := &MockKeyValueEntry{value: []byte("test_value")}
	mockKV.EXPECT().
		Get("test_key").
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
		Get("nonexistent_key").
		Return(nil, nats.ErrKeyNotFound)

	cl := Client{
		kv:      mockKV,
		configs: configs,
	}

	val, err := cl.Get(context.Background(), "nonexistent_key")
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
		Delete("test_key").
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
		Delete("nonexistent_key").
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

	mockJS := NewMockJts(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}

	mockJS.EXPECT().
		AccountInfo().
		Return(&nats.AccountInfo{}, nil)

	cl := Client{
		js:      mockJS,
		configs: configs,
	}

	val, err := cl.HealthCheck(context.Background())
	require.NoError(t, err)

	health := val.(*Health)
	assert.Equal(t, "UP", health.Status)
	assert.Equal(t, configs.Server, health.Details["url"])
	assert.Equal(t, configs.Bucket, health.Details["bucket"])
}

func Test_ClientHealthCheckFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJS := NewMockJts(ctrl)

	configs := &Configs{
		Server: "nats://localhost:4222",
		Bucket: "test_bucket",
	}

	mockJS.EXPECT().
		AccountInfo().
		Return(nil, errConnectionFailed)

	cl := Client{
		js:      mockJS,
		configs: configs,
	}

	val, err := cl.HealthCheck(context.Background())
	require.Error(t, err)
	require.Equal(t, errStatusDown, err)

	health := val.(*Health)
	assert.Equal(t, "DOWN", health.Status)
	assert.Equal(t, configs.Server, health.Details["url"])
	assert.Equal(t, configs.Bucket, health.Details["bucket"])
}
