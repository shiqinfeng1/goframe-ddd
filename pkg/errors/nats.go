package errors

import "github.com/gogf/gf/v2/errors/gcode"

var (
	ErrNatsConnectFail       = codeFunc(gcode.New(CodeNatsConnectFail, "", nil), "nats.connectFail")
	ErrNatsGetStreamInfoFail = codeFunc(gcode.New(CodeNatsGetStreamInfoFail, "", nil), "nats.getStreamInfoFail")
	ErrNatsGetStreamListFail = codeFunc(gcode.New(CodeNatsGetStreamListFail, "", nil), "nats.getStreamListFail")
	ErrNatsDeleteStreamFail  = codeFunc(gcode.New(CodeNatsDeleteStreamFail, "", nil), "nats.getDeleteStreamFail")
)
