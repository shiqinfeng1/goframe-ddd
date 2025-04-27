package server

import "github.com/google/wire"

var WireProviderSet = wire.NewSet(NewHttpServer)
