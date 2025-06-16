package natsclient

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/store"
)

func (c *Client) GetValue(ctx context.Context, bucket, key string) (v string, err error) {
	err = store.SendOperationStats(ctx, store.GET_VAULE, bucket, key, func() error {
		kv := c.kv.Get(bucket)
		if kv == nil {
			return gerror.Newf("bucket not found. bucket=%v", bucket)
		}
		entry, err := kv.(jetstream.KeyValue).Get(ctx, key)
		if err != nil {
			if gerror.Is(err, nats.ErrKeyNotFound) {
				return gerror.Wrapf(nats.ErrKeyNotFound, "bucket=%v, key=%v", bucket, key)
			}
			return gerror.Wrapf(err, "failed to get key-value pair. bucket=%v, key=%v", bucket, key)
		}
		v = string(entry.Value())
		return nil
	})
	return
}

func (c *Client) SetValue(ctx context.Context, bucket, key, value string) (err error) {
	err = store.SendOperationStats(ctx, store.SET_VAULE, bucket, key, func() error {
		kv := c.kv.Get(bucket)
		if kv == nil {
			return gerror.Newf("bucket not found. bucket=%v", bucket)
		}
		_, err := kv.(jetstream.KeyValue).Put(ctx, key, []byte(value))
		if err != nil {
			return gerror.Wrapf(err, "failed to set key-value pair. bucket=%v, key=%v", bucket, key)
		}
		return nil
	})
	return
}

func (c *Client) DeleteValue(ctx context.Context, bucket, key string) error {

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
