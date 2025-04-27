package repositories

import "github.com/google/wire"

var WireProviderSet = wire.NewSet(NewPointmgrRepo)
