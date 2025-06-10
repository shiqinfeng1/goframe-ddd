package main

import (
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/nats-io/nats.go"
	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := gctx.New()
	wg := sync.WaitGroup{}

	// 检查nats消息中间件是否运行
	// 若未运行，则等待
	g.Log().Infof(ctx, "start mgrid server ...")
	for {
		nc, err := nats.Connect(g.Cfg().MustGet(ctx, "nats.serverUrl").String())
		if err != nil {
			g.Log().Warningf(ctx, "wait 2s... connect nats server(%v) fail:%v", g.Cfg().MustGet(ctx, "nats.serverUrl").String(), err)
			time.Sleep(2 * time.Second)
			continue
		}
		if nc.IsConnected() {
			break
		}
	}

	// 初始化http服务
	httpSrv, cleanup1, err := initServer()
	if err != nil {
		g.Log().Panic(ctx, err)
	}
	defer cleanup1()

	// 初始化订阅发布服务
	pubsubSrv, cleanup2, err := initSubAndConsume()
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
		if err := pubsubSrv.Run(ctx); err != nil {
			g.Log().Error(ctx, err)
		}
		g.Log().Infof(ctx, "exit pubsub server ok")
	}()

	// pubsubSrv 需手动关闭
	// http服务本身能监听到信号，无需手动关闭
	signalHandler := func(sig os.Signal) {
		g.Log().Infof(ctx, "signal received:'%v'. gracefully shutting down pubsub service", sig.String())
		if err := pubsubSrv.Stop(ctx); err != nil {
			g.Log().Errorf(ctx, "gracefully shutting down pubsub service fail:%v", err)
		}
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
