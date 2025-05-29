//go:build wireinject

package main

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/service"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/infrastructure/repositories"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server/http"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server/pubsub"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl/dockercmd"
)

// ProvideContext 提供 context.Context 实例
func ProvideCtx() context.Context {
	return gctx.New()
}

var appSrv application.Service

func app(ctx context.Context) (application.Service, error) {
	if appSrv == nil {
		return initApp(ctx)
	}
	return appSrv, nil
}

func initApp(ctx context.Context) (application.Service, error) {
	panic(wire.Build(
		repositories.WireProviderSet,
		service.WireProviderSet,
		application.WireProviderSet,
	))
}

func initServer() (*ghttp.Server, func(), error) {
	panic(wire.Build(
		ProvideCtx,
		dockercmd.WireProviderSet,
		app,
		http.WireProviderSet,
	))
}
func initSubOrConsume() (*pubsub.ControllerV1, func(), error) {
	panic(wire.Build(
		ProvideCtx,
		app,
		pubsub.WireProviderSet,
	))
}
