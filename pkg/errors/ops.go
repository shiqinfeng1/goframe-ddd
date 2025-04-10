package errors

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// code格式： xxyyzz    xx：服务  yy:模块  zz：错误编号
// 30: 运维
//
//	--00: 容器信息查询
//	--01：版本升级
var (
	ErrQueryImageFail   = func(err error) error { return gerror.NewCode(gcode.New(300001, "获取进行信息失败", err)) }
	ErrUpgradeAppFail   = func(err error) error { return gerror.NewCode(gcode.New(300101, "升级应用失败", err)) }
	ErrUpgradeImageFail = func(err error) error { return gerror.NewCode(gcode.New(300102, "升级镜像失败", err)) }
)
