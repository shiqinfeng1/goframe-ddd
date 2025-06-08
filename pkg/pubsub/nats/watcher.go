package natsclient

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type watcher struct {
	logger      pubsub.Logger
	kvWatchers  *gmap.StrAnyMap // map[string]jetstream.KeyWatcher
	objWatchers *gmap.StrAnyMap // map[string]jetstream.ObjectWatcher
	cancel      context.CancelFunc
}
type eventType string

var (
	KV        eventType = "KV"
	OBJECT    eventType = "OBJECT"
	WATCHFAIL eventType = "WATCHFAIL"
)

type ChangeEvent struct {
	Type    eventType // "KV" 或 "OBJECT"
	Bucket  string
	Key     string
	Value   []byte
	Deleted bool
}

func defaultProcessEvent(ctx context.Context, nc *Conn, event ChangeEvent) error {
	eventBytes, _ := json.Marshal(event)
	if event.Type == KV {
		if err := nc.PubMsg(ctx, "notify.change.kv", eventBytes); err != nil {
			return err
		}
	}
	if event.Type == OBJECT {
		if err := nc.PubMsg(ctx, "notify.change.object", eventBytes); err != nil {
			return err
		}
	}
	return nil
}

func (w *watcher) Stop(ctx context.Context) error {
	w.kvWatchers.Iterator(func(key string, value interface{}) bool {
		if err := value.(jetstream.KeyWatcher).Stop(); err != nil {
			w.logger.Errorf(ctx, "stop kv watcher %v failed: %v", key, err)
		}
		return true
	})
	w.objWatchers.Iterator(func(key string, value interface{}) bool {
		if err := value.(jetstream.ObjectWatcher).Stop(); err != nil {
			w.logger.Errorf(ctx, "stop object watcher %v failed: %v", key, err)
		}
		return true
	})
	w.cancel()
	return nil
}
func (w *watcher) StartWatch(
	ctx context.Context,
	nc *Conn, kvbkts, objbkts []string,
	handler func(context.Context, *Conn, ChangeEvent) error) error {
	if handler == nil {
		handler = defaultProcessEvent
	}
	js, err := nc.JetStream()
	if err != nil {
		return err
	}
	// 事件通道
	events := make(chan ChangeEvent, 1024)
	ctx, w.cancel = context.WithCancel(ctx)

	// 启动KV监听
	for _, bkt := range kvbkts {
		go w.watchKV(ctx, js, bkt, events)
	}
	// 启动对象存储监听
	for _, bkt := range objbkts {
		go w.watchObjects(ctx, js, bkt, events)
	}

	// 事件处理器
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-events:
			if event.Type == WATCHFAIL {
				w.logger.Errorf(ctx, "bucket %v %v:%v", event.Bucket, event.Key, string(event.Value))
				continue
			}
			if err := handler(ctx, nc, event); err != nil {
				w.logger.Errorf(ctx, "watch handle fail:%v", err)
			}
		}
	}
}

func (w *watcher) watchKV(ctx context.Context, js jetstream.JetStream, bucket string, events chan<- ChangeEvent) {
	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		events <- ChangeEvent{
			Type:   WATCHFAIL,
			Bucket: bucket,
			Key:    "getkvfail",
			Value:  []byte(err.Error()),
		}
		return
	}
	watcher, err := kv.WatchAll(ctx)
	if err != nil {
		events <- ChangeEvent{
			Type:   WATCHFAIL,
			Bucket: bucket,
			Key:    "watchkvfail",
			Value:  []byte(err.Error()),
		}
		return
	}
	if notexist := w.kvWatchers.SetIfNotExist(bucket, watcher); !notexist {
		w.logger.Errorf(ctx, "kv watcher %v already exist", bucket)
		return
	}

	for entry := range watcher.Updates() {
		if entry != nil {
			events <- ChangeEvent{
				Type:    KV,
				Bucket:  bucket,
				Key:     entry.Key(),
				Value:   entry.Value(),
				Deleted: entry.Operation() != jetstream.KeyValuePut,
			}
		}
	}
}

func (w *watcher) watchObjects(ctx context.Context, js jetstream.JetStream, bucket string, events chan<- ChangeEvent) {
	objStore, err := js.ObjectStore(ctx, bucket)
	if err != nil {
		events <- ChangeEvent{
			Type:   WATCHFAIL,
			Bucket: bucket,
			Key:    "getobjfail",
			Value:  []byte(err.Error()),
		}
		return
	}
	watcher, err := objStore.Watch(ctx)
	if err != nil {
		events <- ChangeEvent{
			Type:   WATCHFAIL,
			Bucket: bucket,
			Key:    "watchobjfail",
			Value:  []byte(err.Error()),
		}
		return
	}
	if notexist := w.objWatchers.SetIfNotExist(bucket, watcher); !notexist {
		w.logger.Errorf(ctx, "obj watcher %v already exist", bucket)
		return
	}

	for info := range watcher.Updates() {
		if info != nil {
			events <- ChangeEvent{
				Type:    OBJECT,
				Bucket:  bucket,
				Key:     info.Name,
				Deleted: info.Deleted,
			}
		}
	}
}
