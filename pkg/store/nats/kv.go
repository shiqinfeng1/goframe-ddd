package nats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) GetValue(ctx context.Context, key string) (string, error) {
	defer c.sendOperationStats(ctx, time.Now(), "GET_VALUE", key)

	entry, err := c.kv.Get(ctx, key)
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			return "", fmt.Errorf("%w: %s", nats.ErrKeyNotFound, key)
		}
		return "", fmt.Errorf("failed to get key: %w", err)
	}
	return string(entry.Value()), nil
}

func (c *Client) SetValue(ctx context.Context, key, value string) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET_VALUE", key)

	_, err := c.kv.Put(ctx, key, []byte(value))
	if err != nil {
		return fmt.Errorf("failed to set key-value pair: %w", err)
	}
	return nil
}

func (c *Client) DeleteValue(ctx context.Context, key string) error {
	defer c.sendOperationStats(ctx, time.Now(), "DELETE_VALUE", key)

	err := c.kv.Delete(ctx, key)
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			return fmt.Errorf("%w: %s", nats.ErrKeyNotFound, key)
		}
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}

func (c *Client) WatchValue(ctx context.Context, key []string) (jetstream.KeyWatcher, error) {
	defer c.sendOperationStats(ctx, time.Now(), "WATCH_VALUE", gconv.String(key))

	watcher, err := c.kv.WatchFiltered(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to watch key: %w", err)
	}
	return watcher, nil
}
