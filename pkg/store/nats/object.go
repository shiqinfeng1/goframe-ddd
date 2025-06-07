package natsclient

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) GetObject(ctx context.Context, bucket, key string) ([]byte, error) {
	defer c.sendOperationStats(ctx, time.Now(), "GET_OBJECT", key)

	entry, err := c.obj[bucket].GetBytes(ctx, key)
	if err != nil {
		if gerror.Is(err, nats.ErrObjectNotFound) {
			return nil, gerror.Wrapf(nats.ErrObjectNotFound, "key=%v", key)
		}
		return nil, gerror.Wrap(err, "failed to get object")
	}
	return entry, nil
}

// fileDir：文件保存路径
func (c *Client) GetFile(ctx context.Context, bucket, key, fileDir string) error {
	defer c.sendOperationStats(ctx, time.Now(), "GET_FILE", key)

	err := c.obj[bucket].GetFile(ctx, key, fileDir)
	if err != nil {
		if gerror.Is(err, nats.ErrObjectNotFound) {
			return gerror.Wrapf(nats.ErrObjectNotFound, "key=%v", key)
		}
		return gerror.Wrap(err, "failed to get file")
	}

	return nil
}
func (c *Client) SetObject(ctx context.Context, bucket, key string, value []byte) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET_OBJECT", key)

	_, err := c.obj[bucket].PutBytes(ctx, key, []byte(value))
	if err != nil {
		return gerror.Wrap(err, "failed to set object")
	}
	return nil
}
func (c *Client) SetFile(ctx context.Context, bucket, file string) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET_OBJECT", file)

	_, err := c.obj[bucket].PutFile(ctx, file)
	if err != nil {
		return gerror.Wrap(err, "failed to set file")
	}
	return nil
}

func (c *Client) DeleteObject(ctx context.Context, bucket, key string) error {
	defer c.sendOperationStats(ctx, time.Now(), "DELETE_OBJECT", key)

	err := c.obj[bucket].Delete(ctx, key)
	if err != nil {
		if gerror.Is(err, nats.ErrObjectNotFound) {
			return gerror.Wrapf(nats.ErrObjectNotFound, "key=%v", key)
		}
		return gerror.Wrap(err, "failed to delete object")
	}
	return nil
}

func (c *Client) WatchObject(ctx context.Context, bucket string, key []string) (jetstream.ObjectWatcher, error) {
	defer c.sendOperationStats(ctx, time.Now(), "WATCH_OBJECT", gconv.String(key))

	watcher, err := c.obj[bucket].Watch(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to watch key")
	}
	return watcher, nil
}
