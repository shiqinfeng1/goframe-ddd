package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := gctx.New()
	httpSrv := server.NewHttpServer()
	g.Log().Infof(ctx, "start http server ...")
	httpSrv.Run()
	g.Log().Infof(ctx, "exit filexfer server \n---------------------------------------------------------\n")
}
