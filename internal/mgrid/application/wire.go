package application

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/pkg/logging"
)

var WireProviderSet = wire.NewSet(New, ProvideLogger)

func ProvideLogger(ctx context.Context) Logger {
	l := glog.New()
	l.SetConfigWithMap(g.Cfg().MustGet(gctx.New(), "logger").Map())
	l.SetHandlers(logging.LoggingGrayLogHandler)
	l.SetPrefix("[APP]")
	return l
}
