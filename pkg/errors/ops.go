package errors

import "github.com/gogf/gf/v2/errors/gcode"

var (
	ErrGetImageListFail = codeFunc(gcode.New(CodeQueryImageFail, "", nil), "ops.queryImageFail")
	ErrUpgradeAppFail   = codeFunc(gcode.New(CodeUpgradeAppFail, "", nil), "ops.upgradeAppFail")
	ErrUpgradeImageFail = codeFunc(gcode.New(CodeUpgradeImageFail, "", nil), "ops.upgradeImageFail")
)
