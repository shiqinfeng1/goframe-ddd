package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := gctx.New()
	httpSrv, cleanup, err := initServer()
	if err != nil {
		g.Log().Panic(ctx, err)
	}
	defer cleanup()
	g.Log().Infof(ctx, "start http server ...")
	httpSrv.Run()
	g.Log().Infof(ctx, "exit filexfer server \n---------------------------------------------------------\n")
}
