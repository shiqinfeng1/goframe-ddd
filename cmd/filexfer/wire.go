//go:build wireinject

package main

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/infrastructure/repositories"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/infrastructure/repositories/migration"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/server"
)

// ProvideContext 提供 context.Context 实例
func ProvideContext() context.Context {
	return gctx.New()
}

func initServer() (*ghttp.Server, func(), error) {
	panic(wire.Build(
		ProvideContext,
		application.WireProviderSet,
		repositories.WireProviderSet,
		server.WireProviderSet,
		filemgr.WireProviderSet,
		migration.WireProviderSet))
}
