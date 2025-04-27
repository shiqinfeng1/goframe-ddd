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
	httpSrv, cleanup1, err := initServer()
	if err != nil {
		g.Log().Panic(ctx, err)
	}
	defer cleanup1()
	pubsubMgr, cleanup2, err := initSub()
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
			g.Log().Fatalf(ctx, "subscription error : %v", err)
		}
		g.Log().Infof(ctx, "exit nats subscrib ok")
	}()

	// grpc服务需要手动关闭
	// submgr 手动关闭
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
