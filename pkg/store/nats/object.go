package natsclient

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/store"
)

func (c *Client) GetObject(ctx context.Context, bucket, key string) (v []byte, err error) {
	err = store.SendOperationStats(ctx, store.GET_OBJ, bucket, key, func() error {
		obj := c.obj.Get(bucket)
		if obj == nil {
			return gerror.Newf("bucket not found. bucket=%v", bucket)
		}
		entry, err := obj.(jetstream.ObjectStore).GetBytes(ctx, key)
		if err != nil {
			return gerror.Wrapf(err, "failed to get object. bucket=%v, key=%v", bucket, key)
		}
		v = entry
		return nil
	})
	return
}

// fileDir：文件保存路径
func (c *Client) GetFile(ctx context.Context, bucket, key, fileDir string) (err error) {
	err = store.SendOperationStats(ctx, store.GET_FILE, bucket, key, func() error {
		obj := c.obj.Get(bucket)
		if obj == nil {
			return gerror.Newf("bucket not found. bucket=%v", bucket)
		}
		err := obj.(jetstream.ObjectStore).GetFile(ctx, key, fileDir)
		if err != nil {
			return gerror.Wrapf(err, "failed to get file. bucket=%v, key=%v", bucket, key)
		}
		return nil
	})
	return
}
func (c *Client) SetObject(ctx context.Context, bucket, key string, value []byte) (err error) {
	err = store.SendOperationStats(ctx, store.SET_OBJ, bucket, key, func() error {
		obj := c.obj.Get(bucket)
		if obj == nil {
			return gerror.Newf("bucket not found. bucket=%v", bucket)
		}
		_, err := obj.(jetstream.ObjectStore).PutBytes(ctx, key, []byte(value))
		if err != nil {
			return gerror.Wrapf(err, "failed to set object. bucket=%v, key=%v", bucket, key)
		}
		return nil
	})
	return
}
func (c *Client) SetFile(ctx context.Context, bucket, file string) (err error) {
	err = store.SendOperationStats(ctx, store.SET_FILE, bucket, file, func() error {
		obj := c.obj.Get(bucket)
		if obj == nil {
			return gerror.Newf("bucket not found. bucket=%v", bucket)
		}
		_, err := obj.(jetstream.ObjectStore).PutFile(ctx, file)
		if err != nil {
			return gerror.Wrapf(err, "failed to set file. bucket=%v, file=%v", bucket, file)
		}
		return nil
	})
	return
}

func (c *Client) DeleteObject(ctx context.Context, bucket, key string) error {
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
