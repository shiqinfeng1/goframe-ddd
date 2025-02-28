package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/shiqinfeng1/goframe-ddd/internal/server"
	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := gctx.New()
	go func() {
		g.Log().Infof(ctx, "start http server ...")
		server.NewHttpServer().Run()
		g.Log().Infof(ctx, "exit http server ok")
	}()

	go func() {
		g.Log().Infof(ctx, "start grpc server ...")
		server.NewGrpcServer().Run()
		g.Log().Infof(ctx, "exit grpc server ok")
	}()

	gproc.AddSigHandlerShutdown()
	gproc.Listen()
}
