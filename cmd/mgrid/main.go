package main

import (
	"os"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/shiqinfeng1/goframe-ddd/internal/server"
	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := gctx.New()
	wg := sync.WaitGroup{}
	httpSrv := server.NewHttpServer()
	grpcSrv := server.NewGrpcServer()
	subMgr := server.NewSubscriptions()

	// grpc服务需要手动关闭
	signalHandler := func(sig os.Signal) {
		g.Log().Infof(ctx, "signal received: %v, gracefully shutting down grpc server", sig.String())
		grpcSrv.Stop()
		subMgr.Stop()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Log().Infof(ctx, "start http server ...")
		httpSrv.Run()
		g.Log().Infof(ctx, "exit http server ok")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Log().Infof(ctx, "start grpc server ...")
		grpcSrv.Run()
		g.Log().Infof(ctx, "exit grpc server ok")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Log().Infof(ctx, "start nats subscrib ...")
		if err := subMgr.Run(ctx); err != nil {
			g.Log().Fatal(ctx, "subscription error : %v", err)
		}
		g.Log().Infof(ctx, "exit nats subscrib ok")
	}()

	// 监听系统中断信号
	gproc.AddSigHandlerShutdown(
		signalHandler,
	)
	gproc.Listen()

	wg.Wait()
	g.Log().Infof(ctx, "exit all ok")
}
