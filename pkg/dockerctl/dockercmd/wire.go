package dockercmd

import "github.com/google/wire"

var WireProviderSet = wire.NewSet(New)
