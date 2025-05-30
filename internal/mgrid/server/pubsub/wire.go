package pubsub

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	"github.com/shiqinfeng1/goframe-ddd/pkg/logging"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

var WireProviderSet = wire.NewSet(NewV1, ProvideLogger)
var WireProviderNatsFactory = wire.NewSet(ProvideLogger, ProvideNatsServerAddr, ProvideConnFactory)

func ProvideLogger() server.Logger {
	l := g.Log()
	l.SetPrefix("pubsub")
	l.SetAsync(true)
	l.SetHandlers(logging.LoggingJsonHandler)
	return l
}
func ProvideConnFactory(logger server.Logger, natsAddr string) nats.ConnFactory {
	return nats.NewFactory(logger, natsAddr, nil)
}

func ProvideNatsServerAddr(ctx context.Context) string {
	return g.Cfg().MustGet(ctx, "nats.serverAddr").String()
}
