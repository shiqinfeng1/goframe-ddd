package repositories

import "github.com/google/wire"

var WireProviderSet = wire.NewSet(NewPointdataRepo, NewTokenRepo, NewUserRepo)
