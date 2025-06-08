package natsclient

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) GetObject(ctx context.Context, bucket, key string) ([]byte, error) {
	defer c.sendOperationStats(ctx, time.Now(), "GET_OBJECT", key)

	obj := c.obj.Get(bucket)
	if obj == nil {
		return nil, gerror.Newf("bucket not found. bucket=%v", bucket)
	}
	entry, err := obj.(jetstream.ObjectStore).GetBytes(ctx, key)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to get object. bucket=%v, key=%v", bucket, key)
	}
	return entry, nil
}

// fileDir：文件保存路径
func (c *Client) GetFile(ctx context.Context, bucket, key, fileDir string) error {
	defer c.sendOperationStats(ctx, time.Now(), "GET_FILE", key)

	obj := c.obj.Get(bucket)
	if obj == nil {
		return gerror.Newf("bucket not found. bucket=%v", bucket)
	}
	err := obj.(jetstream.ObjectStore).GetFile(ctx, key, fileDir)
	if err != nil {
		return gerror.Wrapf(err, "failed to get file. bucket=%v, key=%v", bucket, key)
	}

	return nil
}
func (c *Client) SetObject(ctx context.Context, bucket, key string, value []byte) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET_OBJECT", key)
	obj := c.obj.Get(bucket)
	if obj == nil {
		return gerror.Newf("bucket not found. bucket=%v", bucket)
	}
	_, err := obj.(jetstream.ObjectStore).PutBytes(ctx, key, []byte(value))
	if err != nil {
		return gerror.Wrapf(err, "failed to set object. bucket=%v, key=%v", bucket, key)
	}
	return nil
}
func (c *Client) SetFile(ctx context.Context, bucket, file string) error {
	defer c.sendOperationStats(ctx, time.Now(), "SET_OBJECT", file)

	obj := c.obj.Get(bucket)
	if obj == nil {
		return gerror.Newf("bucket not found. bucket=%v", bucket)
	}
	_, err := obj.(jetstream.ObjectStore).PutFile(ctx, file)
	if err != nil {
		return gerror.Wrapf(err, "failed to set file. bucket=%v, file=%v", bucket, file)
	}
	return nil
}

func (c *Client) DeleteObject(ctx context.Context, bucket, key string) error {
	defer c.sendOperationStats(ctx, time.Now(), "DELETE_OBJECT", key)

	obj := c.obj.Get(bucket)
	if obj == nil {
		return gerror.Newf("bucket not found. bucket=%v", bucket)
	}
	err := obj.(jetstream.ObjectStore).Delete(ctx, key)
	if err != nil {
		return gerror.Wrapf(err, "failed to delete object. bucket=%v, key=%v", bucket, key)
	}
	return nil
}
