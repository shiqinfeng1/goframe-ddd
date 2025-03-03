package errors

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// code格式： xxyyzz    xx：服务  yy:模块  zz：错误编号
// 10：文件管理服务
// --00：发送模块
//
// 20：
// --00：
var (
	ErrInvalidFile = func(f string) error { return gerror.NewCode(gcode.New(100001, "无效的文件:"+f, nil)) }
)
