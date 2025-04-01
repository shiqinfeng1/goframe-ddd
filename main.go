package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/grand"
	"golang.org/x/time/rate"
)

func main() {
	ctx := gctx.New()

	msg := []byte(grand.Letters(128))

	cancel := make(chan struct{})
	oneDay := 24 * time.Hour
	time.AfterFunc(oneDay, func() {
		close(cancel)
	})
	// 限速每秒发送50个
	limiter := rate.NewLimiter(rate.Limit(1), 1)
	for {
		select {
		case <-cancel:
			return
		default:
			limiter.Wait(ctx)
			fmt.Printf("%v pub msg:%s\n", time.Now(), msg)
		}
	}
}
