package main

import (
	"os"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := gctx.New()
	wg := sync.WaitGroup{}

	// 初始化http服务
	httpSrv, cleanup1, err := initServer()
	if err != nil {
		g.Log().Panic(ctx, err)
	}
	defer cleanup1()

	// 初始化订阅发布服务
	pubsubMgr, cleanup2, err := initSubOrConsume()
	if err != nil {
		g.Log().Panic(ctx, err)
	}
	defer cleanup2()

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpSrv.Run()
		g.Log().Infof(ctx, "exit http server ok")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := pubsubMgr.Run(ctx); err != nil {
			g.Log().Fatalf(ctx, "pubsub run error : %v", err)
		}
		g.Log().Infof(ctx, "exit nats subscrib ok")
	}()

	// pubsubMgr 需手动关闭
	// http服务本身能监听到信号，无需手动关闭
	signalHandler := func(sig os.Signal) {
		g.Log().Infof(ctx, "signal received: @@@@ '%v' @@@@, gracefully shutting down pubsub service", sig.String())
		pubsubMgr.Stop(ctx)
		g.Log().Infof(ctx, "gracefully shutting down pubsub service ok")
	}
	// 监听系统中断信号
	gproc.AddSigHandlerShutdown(
		signalHandler,
	)
	gproc.Listen()

	wg.Wait()
	g.Log().Infof(ctx, "exit all ok\n---------------------------------------------------------\n")
}
