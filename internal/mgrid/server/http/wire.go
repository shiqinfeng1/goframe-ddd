package http

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	"github.com/shiqinfeng1/goframe-ddd/pkg/logging"
)

var WireProviderSet = wire.NewSet(NewServer, ProvideLogger)

func ProvideLogger() server.Logger {
	l := g.Log()
	l.SetPrefix("http-serevr")
	l.SetAsync(true)
	l.SetHandlers(logging.LoggingJsonHandler)
	return l
}
