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
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server/pubsub"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl/dockercmd"
)

// ProvideContext 提供 context.Context 实例
func ProvideContext() context.Context {
	return gctx.New()
}

func initServer() (*ghttp.Server, func(), error) {
	panic(wire.Build(
		ProvideContext,
		service.WireProviderSet,
		application.WireProviderSet,
		repositories.WireProviderSet,
		server.WireHttpProviderSet,
		dockercmd.WireProviderSet))
}
func initSub() (*pubsub.ControllerV1, func(), error) {
	panic(wire.Build(
		ProvideContext,
		service.WireProviderSet,
		repositories.WireProviderSet,
		application.WireProviderSet,
		server.WireSubProviderSet))
}
