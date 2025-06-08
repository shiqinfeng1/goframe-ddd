package natsclient

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) GetValue(ctx context.Context, bucket, key string) (string, error) {
	defer c.sendOperationStats(ctx, time.Now(), "GET_VALUE", key)

	kv := c.kv.Get(bucket)
	if kv == nil {
		return "", gerror.Newf("bucket not found. bucket=%v", bucket)
	}
	entry, err := kv.(jetstream.KeyValue).Get(ctx, key)
	if err != nil {
		if gerror.Is(err, nats.ErrKeyNotFound) {
			return "", gerror.Wrapf(nats.ErrKeyNotFound, "bucket=%v, key=%v", bucket, key)
		}
		return "", gerror.Wrapf(err, "failed to get key-value pair. bucket=%v, key=%v", bucket, key)
	}
	return string(entry.Value()), nil
}

func (c *Client) SetValue(ctx context.Context, bucket, key, value string) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET_VALUE", key)

	kv := c.kv.Get(bucket)
	if kv == nil {
		return gerror.Newf("bucket not found. bucket=%v", bucket)
	}
	_, err := kv.(jetstream.KeyValue).Put(ctx, key, []byte(value))
	if err != nil {
		return gerror.Wrapf(err, "failed to set key-value pair. bucket=%v, key=%v", bucket, key)
	}
	return nil
}

func (c *Client) DeleteValue(ctx context.Context, bucket, key string) error {
	defer c.sendOperationStats(ctx, time.Now(), "DELETE_VALUE", key)

	kv := c.kv.Get(bucket)
	if kv == nil {
		return gerror.Newf("bucket not found. bucket=%v", bucket)
	}
	err := kv.(jetstream.KeyValue).Delete(ctx, key)
	if err != nil {
		return gerror.Wrapf(err, "failed to delete key-value pair. bucket=%v, key=%v", bucket, key)
	}
	return nil
}
