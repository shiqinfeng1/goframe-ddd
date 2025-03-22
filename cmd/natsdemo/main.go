package main

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func expectOk(ctx context.Context, err error) {
	if err != nil {
		g.Log().Fatalf(ctx, "Unexpected error: %v", err)
	}
}

func expectErr(ctx context.Context, err error, expected ...error) {
	if err == nil {
		g.Log().Fatalf(ctx, "Expected error but got none")
	}
	if len(expected) == 0 {
		return
	}
	for _, e := range expected {
		if errors.Is(err, e) {
			return
		}
	}
	g.Log().Fatalf(ctx, "Expected one of %+v, got '%v'", expected, err)
}
func main() {
	ctx := gctx.New()
	nc, err := nats.Connect("http://localhost:4222")
	if err != nil {
		g.Log().Fatalf(ctx, "Unexpected error: %v", err)
	}
	defer nc.Close()
	g.Log().Info(ctx, "connect nats server ok: http://localhost:4222")
	js, err := jetstream.New(nc)
	if err != nil {
		g.Log().Fatalf(ctx, "Unexpected error getting JetStream context: %v", err)
	}
	g.Log().Info(ctx, "new jetstream ok")

	kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "TEST1", History: 5, TTL: time.Hour})
	expectOk(ctx, err)
	g.Log().Info(ctx, "创建一个kv桶：TEST，保留5个版本，过期时间：1小时")
	if kv.Bucket() != "TEST1" {
		g.Log().Fatalf(ctx, "Expected bucket name1 to be %q, got %q", "TEST1", kv.Bucket())
	}

	// Simple Put
	r, err := kv.Put(ctx, "name1", []byte("derek"))
	expectOk(ctx, err)
	if r != 1 {
		g.Log().Fatalf(ctx, "Expected 1 for the revision, got %d", r)
	}
	g.Log().Info(ctx, "通过put方法保存 k=name1")
	g.Log().Info(ctx, "put k=name1 v=derek ok")
	// Simple Get
	e, err := kv.Get(ctx, "name1")
	expectOk(ctx, err)
	if string(e.Value()) != "derek" {
		g.Log().Fatalf(ctx, "Got wrong value: %q vs %q", e.Value(), "derek")
	}
	if e.Revision() != 1 {
		g.Log().Fatalf(ctx, "Expected 1 for the revision, got %d", e.Revision())
	}
	g.Log().Infof(ctx, "get k=name1 v=%s revision=%v ok", e.Value(), e.Revision())

	// Delete
	err = kv.Delete(ctx, "name1")
	expectOk(ctx, err)
	g.Log().Info(ctx, "delete k=name1 v=%v  ok")

	_, err = kv.Get(ctx, "name1")
	expectErr(ctx, err, jetstream.ErrKeyNotFound)

	r, err = kv.Create(ctx, "name1", []byte("derek"))
	expectOk(ctx, err)
	if r != 3 {
		g.Log().Fatalf(ctx, "Expected 3 for the revision, got %d", r)
	}
	g.Log().Info(ctx, "通过create方法保存 k=name1")
	g.Log().Infof(ctx, "create k=name1 v=%s revision=%v ok", e.Value(), e.Revision())

	err = kv.Delete(ctx, "name1", jetstream.LastRevision(4))
	expectErr(ctx, err)
	err = kv.Delete(ctx, "name1", jetstream.LastRevision(3))
	expectOk(ctx, err)

	// Conditional Updates.
	r, err = kv.Update(ctx, "name1", []byte("rip"), 4)
	expectOk(ctx, err)
	g.Log().Infof(ctx, "update k=name1 v=%s revision=%v ok", e.Value(), e.Revision())
	_, err = kv.Update(ctx, "name1", []byte("ik"), 3)
	expectErr(ctx, err)
	_, err = kv.Update(ctx, "name1", []byte("ik"), r)
	expectOk(ctx, err)
	r, err = kv.Create(ctx, "age1", []byte("22"))
	expectOk(ctx, err)
	g.Log().Infof(ctx, "create k=age1 v=22 revision=%v ok", e.Revision())
	_, err = kv.Update(ctx, "age1", []byte("33"), r)
	expectOk(ctx, err)
	g.Log().Infof(ctx, "update k=age1 v=33 revision=%v ok", e.Revision())

	// Status
	status, err := kv.Status(ctx)
	expectOk(ctx, err)
	if status.History() != 5 {
		g.Log().Fatalf(ctx, "expected history of 5 got %d", status.History())
	}
	if status.Bucket() != "TEST1" {
		g.Log().Fatalf(ctx, "expected bucket TEST1 got %v", status.Bucket())
	}
	if status.TTL() != time.Hour {
		g.Log().Fatalf(ctx, "expected 1 hour TTL got %v", status.TTL())
	}
	if status.Values() != 7 {
		g.Log().Fatalf(ctx, "expected 7 values got %d", status.Values())
	}
	if status.BackingStore() != "JetStream" {
		g.Log().Fatalf(ctx, "invalid backing store kind %s", status.BackingStore())
	}

	kvs := status.(*jetstream.KeyValueBucketStatus)
	si := kvs.StreamInfo()
	if si == nil {
		g.Log().Fatalf(ctx, "StreamInfo not received")
	}
	g.Log().Infof(ctx, "桶：TEST status=%v ok", status)
}
