package service

import "github.com/google/wire"

var WireProviderSet = wire.NewSet(NewPointdataService, NewJetstreamService, NewAuthService)
