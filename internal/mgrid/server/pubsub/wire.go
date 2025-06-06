package pubsub

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	"github.com/shiqinfeng1/goframe-ddd/pkg/logging"
	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

var WireProviderSet = wire.NewSet(NewV1, ProvideLogger)
var WireProviderNatsFactory = wire.NewSet(ProvideLogger, ProvideConnFactory)

func ProvideLogger() server.Logger {
	l := g.Log()
	l.SetPrefix("eventServer")
	l.SetAsync(true)
	l.SetHandlers(logging.JsonHandler)
	return l
}
func ProvideConnFactory(logger server.Logger) natsclient.Factory {
	return natsclient.NewFactory(logger, nil)
}
