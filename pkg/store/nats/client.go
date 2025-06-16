package natsclient

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/store"
)

type Config struct {
	Server     string   `json:"server" yaml:"server"`
	KVBuckets  []string `json:"kvBuckets" yaml:"kvBuckets"`
	ObjBuckets []string `json:"objBuckets" yaml:"objBuckets"`
}

type Client struct {
	logger store.Logger
	conn   *nats.Conn
	js     jetstream.JetStream
	kv     *gmap.StrAnyMap //map[string]jetstream.KeyValue
	obj    *gmap.StrAnyMap //map[string]jetstream.ObjectStore
	cfg    *Config
}

// New creates a new NATS-KV client with the provided configuration.
func New(logger store.Logger) *Client {
	c := &Client{
		logger: logger,
	}
	ctx := gctx.New()
	if err := g.Cfg().MustGet(ctx, "store").Scan(&c.cfg); err != nil {
		logger.Fatalf(ctx, "get nats config fail:%v", err)
	}
	return c
}
