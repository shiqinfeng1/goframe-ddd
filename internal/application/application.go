package application

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters"
	"github.com/shiqinfeng1/goframe-ddd/internal/application/command"
	"github.com/shiqinfeng1/goframe-ddd/internal/application/query"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/pkg/stream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/stream/transport"
)

type Application struct {
	Commands *command.Handler
	Queries  *query.Handler
}

// New 一个DDD的应用层，支持CQRS，
func New(ctx context.Context) *Application {
	// 实例化一个文件管理的数据仓库
	repo := adapters.NewFilemgrRepo()

	// 实例化一个文件传输服务
	maxTasks := g.Cfg().MustGet(ctx, "filemgr.maxTasks").Int()
	fileTransferService := filemgr.NewfileTransferService(maxTasks)

	// 实例化一个流通道管理服务，流通道支持2种传输层：tcp和kcp
	var streamService *stream.Stream
	isCloud := g.Cfg().MustGet(ctx, "filemgr.isCloud").Bool()
	addr := g.Cfg().MustGet(ctx, "filemgr.addr").String()
	transType := g.Cfg().MustGet(ctx, "filemgr.transport").String()
	switch transType {
	case "kcp":
		streamService = stream.New(ctx, transport.NewKcpTransport())
	case "tcp":
		streamService = stream.New(ctx, transport.NewTcpTransport())
	default:
		g.Log().Fatalf(ctx, "config filemgr.transport is invalid:%v", transType)
	}
	if isCloud {
		streamService.StartupServer(ctx, addr, filemgr.StreamRecvHandler)
	} else {
		streamService.StartupClient(ctx, addr)
	}

	return &Application{
		Commands: command.NewHandler(repo, fileTransferService, streamService),
		Queries:  query.NewHandler(repo),
	}
}
