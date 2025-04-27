package service

import "github.com/google/wire"

var WireProviderSet = wire.NewSet(NewPointDataSetService, NeJetStreamService)
