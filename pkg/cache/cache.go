package cache

import "github.com/gogf/gf/v2/os/gcache"

var inst *gcache.Cache

func init() {
	inst = gcache.NewWithAdapter(gcache.NewAdapterMemory())
}

func Memory() *gcache.Cache {
	return inst
}
