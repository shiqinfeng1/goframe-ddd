package migration

import "github.com/google/wire"

var WireProviderSet = wire.NewSet(NewEntClient)
