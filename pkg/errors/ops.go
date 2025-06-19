package errors

import "github.com/gogf/gf/v2/errors/gcode"

var (
	ErrGetImageListFail = codeFunc(gcode.New(CodeQueryImageFail, "", nil), "opsQueryImageFail")
	ErrUpgradeAppFail   = codeFunc(gcode.New(CodeUpgradeAppFail, "", nil), "opsUpgradeAppFail")
	ErrUpgradeImageFail = codeFunc(gcode.New(CodeUpgradeImageFail, "", nil), "opsUpgradeImageFail")
)
