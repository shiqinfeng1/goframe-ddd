package application

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/wire"
	"github.com/shiqinfeng1/goframe-ddd/pkg/logging"
)

var WireProviderSet = wire.NewSet(New, ProvideLogger)

func ProvideLogger() Logger {
	l := g.Log()
	l.SetPrefix("app")
	l.SetAsync(true)
	l.SetHandlers(logging.JsonHandler)
	return l
}
