package errors

import "github.com/gogf/gf/v2/errors/gcode"

var (
	ErrAuthRefreshTokenFail = codeFunc(gcode.New(CodeAuthRefreshTokenFail, "", nil), "auth.refreshTokenFail")
	ErrSendVCodeFail        = codeFunc(gcode.New(CodeSendVCodeFail, "", nil), "auth.sendVCodeFail")
	ErrRestPwdFail          = codeFunc(gcode.New(CodeResetPwdFail, "", nil), "auth.resetPwdFail")
	ErrRegisterUserFail     = codeFunc(gcode.New(CodeRegisterUserFail, "", nil), "auth.registerUserFail")
	ErrLoginFail            = codeFunc(gcode.New(CodeLoginFail, "", nil), "auth.loginFail")
	ErrLogoutFail           = codeFunc(gcode.New(CodeLogoutFail, "", nil), "auth.logoutFail")
)
