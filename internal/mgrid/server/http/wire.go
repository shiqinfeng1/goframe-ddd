package http

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	"github.com/shiqinfeng1/goframe-ddd/pkg/logging"
)

var WireProviderSet = wire.NewSet(NewServer, ProvideLogger)

func ProvideLogger(ctx context.Context) server.Logger {
	l := glog.New()
	l.SetConfigWithMap(g.Cfg().MustGet(ctx, "logger").Map())
	l.SetHandlers(logging.LoggingGrayLogHandler)
	l.SetPrefix("[HTTP-SERVER]")
	return l
}
