package natsclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) GetObject(ctx context.Context, key string) ([]byte, error) {
	defer c.sendOperationStats(ctx, time.Now(), "GET_OBJECT", key)

	entry, err := c.obj.GetBytes(ctx, key)
	if err != nil {
		if errors.Is(err, nats.ErrObjectNotFound) {
			return nil, fmt.Errorf("%w: %s", nats.ErrObjectNotFound, key)
		}
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return entry, nil
}

// fileDir：文件保存路径
func (c *Client) GetFile(ctx context.Context, key, fileDir string) error {
	defer c.sendOperationStats(ctx, time.Now(), "GET_FILE", key)

	err := c.obj.GetFile(ctx, key, fileDir)
	if err != nil {
		if errors.Is(err, nats.ErrObjectNotFound) {
			return fmt.Errorf("%w: %s", nats.ErrObjectNotFound, key)
		}
		return fmt.Errorf("failed to get file: %w", err)
	}

	return nil
}
func (c *Client) SetObject(ctx context.Context, key string, value []byte) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET_OBJECT", key)

	_, err := c.obj.PutBytes(ctx, key, []byte(value))
	if err != nil {
		return fmt.Errorf("failed to set object: %w", err)
	}
	return nil
}
func (c *Client) SetFile(ctx context.Context, file string) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET_OBJECT", file)

	_, err := c.obj.PutFile(ctx, file)
	if err != nil {
		return fmt.Errorf("failed to set file: %w", err)
	}
	return nil
}

func (c *Client) DeleteObject(ctx context.Context, key string) error {
	defer c.sendOperationStats(ctx, time.Now(), "DELETE_OBJECT", key)

	err := c.obj.Delete(ctx, key)
	if err != nil {
		if errors.Is(err, nats.ErrObjectNotFound) {
			return fmt.Errorf("%w: %s", nats.ErrObjectNotFound, key)
		}
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

func (c *Client) WatchObject(ctx context.Context, key []string) (jetstream.ObjectWatcher, error) {
	defer c.sendOperationStats(ctx, time.Now(), "WATCH_OBJECT", gconv.String(key))

	watcher, err := c.obj.Watch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to watch key: %w", err)
	}
	return watcher, nil
}
