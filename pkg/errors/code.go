package errors

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/shiqinfeng1/goframe-ddd/pkg/locale"
)

func codeFunc(c gcode.Code, token string) func(string) error {
	return func(l string) error {
		if v, ok := locale.Lang[l]; ok {
			return gerror.NewCode(c, v.T(ctx, token))
		}
		return gerror.NewCode(c, locale.Lang["en"].T(ctx, token))
	}
}
func codefFunc(c gcode.Code, token string) func(string, ...any) error {
	return func(l string, values ...any) error {
		if v, ok := locale.Lang[l]; ok {
			return gerror.NewCode(c, v.Tf(ctx, token, values...))
		}
		return gerror.NewCode(c, locale.Lang["en"].Tf(ctx, token, values...))
	}
}

var ctx = context.Background()

const (
	// code格式： xxyyzz    xx：服务  yy:模块  zz：错误编号
	// 10：文件传输服务
	//	--00：发送模块
	//	--01: 文件读写模块
	CodeNotAbsFilePath = 100001
	CodeEmptyDir       = 100002
	CodeInvalidFiles   = 100003
	CodeInvalidNodeId  = 100004
	CodeOver4GSize     = 100101
	// 11：消息队列
	//	--00：nats管理
	CodeNatsConnectFail       = 110001
	CodeNatsGetStreamInfoFail = 110002
	CodeNatsGetStreamListFail = 110003
	CodeNatsDeleteStreamFail  = 110003
	CodeNatsNotFooundStream   = 110004
	// 12: 运维
	//	--00: 容器信息查询
	//	--01：版本升级
	CodeQueryImageFail   = 120001
	CodeUpgradeAppFail   = 120101
	CodeUpgradeImageFail = 120102

	// 13: 认证
	//	--00: 用户管理
	//	--01: token管理
	CodeSendVCodeFail              = 130001
	CodeResetPwdFail               = 130002
	CodeRegisterUserFail           = 130003
	CodeLoginFail                  = 130004
	CodeLogoutFail                 = 130005
	CodeUserExisted                = 130006
	CodeUserIsLockdBefore          = 130007
	CodeUserIsLockdTooManyAttempts = 130008
	CodeUserVerifyAttemptsRemain   = 130009
	CodeAuthRefreshTokenFail       = 130101
)
