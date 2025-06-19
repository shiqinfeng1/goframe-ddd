package errors

import "github.com/gogf/gf/v2/errors/gcode"

var (
	ErrNatsConnectFail       = codeFunc(gcode.New(CodeNatsConnectFail, "", nil), "natsConnectFail")
	ErrNatsGetStreamInfoFail = codeFunc(gcode.New(CodeNatsGetStreamInfoFail, "", nil), "natsGetStreamInfoFail")
	ErrNatsGetStreamListFail = codeFunc(gcode.New(CodeNatsGetStreamListFail, "", nil), "natsGetStreamListFail")
	ErrNatsDeleteStreamFail  = codeFunc(gcode.New(CodeNatsDeleteStreamFail, "", nil), "natsGetDeleteStreamFail")
)
