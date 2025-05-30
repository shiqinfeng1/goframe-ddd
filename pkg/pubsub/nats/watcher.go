package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
)

type watcher struct {
	logger      pubsub.Logger
	kvWatchers  map[string]jetstream.KeyWatcher
	objWatchers map[string]jetstream.ObjectWatcher
	cancel      context.CancelFunc
}
type eventType string

var (
	KV     eventType = "KV"
	OBJECT eventType = "OBJECT"
)

type ChangeEvent struct {
	Type    eventType // "KV" 或 "OBJECT"
	Bucket  string
	Key     string
	Value   []byte
	Deleted bool
}

func defaultProcessEvent(ctx context.Context, nc *Conn, event ChangeEvent) {
	eventBytes, _ := json.Marshal(event)
	if event.Type == KV {
		nc.PubMsg(ctx, "notify.change.kv", eventBytes)
	}
	if event.Type == OBJECT {
		nc.PubMsg(ctx, "notify.change.object", eventBytes)
	}
}

func (w *watcher) Stop() error {
	for _, w := range w.kvWatchers {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	for _, w := range w.objWatchers {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	w.cancel()
	return nil
}
func (w *watcher) StartWatch(ctx context.Context, nc *Conn, kvbkts, objbkts []string, handler func(context.Context, *Conn, ChangeEvent)) error {
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
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-events:
				handler(ctx, nc, event)
			}
		}
	}()
	return nil
}

func (w *watcher) watchKV(ctx context.Context, js jetstream.JetStream, bucket string, events chan<- ChangeEvent) {
	kv, _ := js.KeyValue(ctx, bucket)
	watcher, _ := kv.WatchAll(ctx)
	w.kvWatchers[bucket] = watcher

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
	objStore, _ := js.ObjectStore(ctx, bucket)
	watcher, _ := objStore.Watch(ctx)
	w.objWatchers[bucket] = watcher

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
