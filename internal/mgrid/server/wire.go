package server

import "github.com/google/wire"

var WireHttpProviderSet = wire.NewSet(NewHttpServer)
var WireSubProviderSet = wire.NewSet(NewSubscriptions)
