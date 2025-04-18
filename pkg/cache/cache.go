package cache

import "github.com/gogf/gf/v2/os/gcache"

type Cache struct {
	kv *gcache.Cache
}

var DefaultBackend Cache

func init() {
	DefaultBackend = Cache{
		kv: gcache.NewWithAdapter(gcache.NewAdapterMemory()),
	}
}

func KV() *gcache.Cache {
	return DefaultBackend.kv
}
