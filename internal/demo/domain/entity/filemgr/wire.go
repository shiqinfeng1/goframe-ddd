package filemgr

import "github.com/google/wire"

var WireProviderSet = wire.NewSet(NewFileTransferService)
