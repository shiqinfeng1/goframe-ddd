package pubsub

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	"github.com/shiqinfeng1/goframe-ddd/pkg/logging"
	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

var WireProviderSet = wire.NewSet(NewV1, ProvideLogger)
var WireProviderNatsFactory = wire.NewSet(ProvideLogger, ProvideConnFactory)

func ProvideLogger(ctx context.Context) server.Logger {
	l := glog.New()
	l.SetConfigWithMap(g.Cfg().MustGet(ctx, "logger").Map())
	l.SetHandlers(logging.LoggingGrayLogHandler)
	l.SetPrefix("[EVENT-SERVER]")
	return l
}
func ProvideConnFactory(logger server.Logger) natsclient.Factory {
	return natsclient.NewFactory(logger, nil)
}
