package errors

import "github.com/gogf/gf/v2/errors/gcode"

var (
	ErrAuthRefreshTokenFail       = codeFunc(gcode.New(CodeAuthRefreshTokenFail, "", nil), "authRefreshTokenFail")
	ErrSendVCodeFail              = codeFunc(gcode.New(CodeSendVCodeFail, "", nil), "authSendVCodeFail")
	ErrRestPwdFail                = codeFunc(gcode.New(CodeResetPwdFail, "", nil), "authResetPwdFail")
	ErrRegisterUserFail           = codeFunc(gcode.New(CodeRegisterUserFail, "", nil), "authRegisterUserFail")
	ErrUserExisted                = codeFunc(gcode.New(CodeUserExisted, "", nil), "authUserExisted")
	ErrLoginFail                  = codeFunc(gcode.New(CodeLoginFail, "", nil), "authLoginFail")
	ErrUserIsLockdBefore          = codefFunc(gcode.New(CodeUserIsLockdBefore, "", nil), "authUserIsLockdBefore")
	ErrUserIsLockdTooManyAttempts = codeFunc(gcode.New(CodeUserIsLockdTooManyAttempts, "", nil), "authUserIsLockdTooManyAttempts")
	ErrUserVerifyAttemptsRemain   = codefFunc(gcode.New(CodeUserVerifyAttemptsRemain, "", nil), "authUserVerifyAttemptsRemain")
	ErrLogoutFail                 = codeFunc(gcode.New(CodeLogoutFail, "", nil), "authLogoutFail")
)
