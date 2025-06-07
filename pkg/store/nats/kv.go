package natsclient

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) GetValue(ctx context.Context, bucket, key string) (string, error) {
	defer c.sendOperationStats(ctx, time.Now(), "GET_VALUE", key)

	entry, err := c.kv[bucket].Get(ctx, key)
	if err != nil {
		if gerror.Is(err, nats.ErrKeyNotFound) {
			return "", gerror.Wrapf(nats.ErrKeyNotFound, "key=%v", key)
		}
		return "", gerror.Wrap(err, "failed to get key")
	}
	return string(entry.Value()), nil
}

func (c *Client) SetValue(ctx context.Context, bucket, key, value string) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET_VALUE", key)

	_, err := c.kv[bucket].Put(ctx, key, []byte(value))
	if err != nil {
		return gerror.Wrapf(err, "failed to set key-value pair")
	}
	return nil
}

func (c *Client) DeleteValue(ctx context.Context, bucket, key string) error {
	defer c.sendOperationStats(ctx, time.Now(), "DELETE_VALUE", key)

	err := c.kv[bucket].Delete(ctx, key)
	if err != nil {
		if gerror.Is(err, nats.ErrKeyNotFound) {
			return gerror.Wrapf(nats.ErrKeyNotFound, "key=%v", key)
		}
		return gerror.Wrap(err, "failed to delete key")
	}
	return nil
}

func (c *Client) WatchValue(ctx context.Context, bucket string, key []string) (jetstream.KeyWatcher, error) {
	defer c.sendOperationStats(ctx, time.Now(), "WATCH_VALUE", gconv.String(key))

	watcher, err := c.kv[bucket].WatchFiltered(ctx, key)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to watch key")
	}
	return watcher, nil
}
