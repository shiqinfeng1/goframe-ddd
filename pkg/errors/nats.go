package errors

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// code格式： xxyyzz    xx：服务  yy:模块  zz：错误编号
// 10：文件传输服务
// --00：发送模块
// --01: 文件读写模块
// 20：消息队列
// --00：nats管理
var (
	ErrNatsConnectFail      = func(err error) error { return gerror.NewCode(gcode.New(200001, "连接nats失败", err)) }
	ErrNatsStreamFail       = func(err error) error { return gerror.NewCode(gcode.New(200002, "查询stream失败", err)) }
	ErrNatsDeleteStreamFail = func(err error) error { return gerror.NewCode(gcode.New(200003, "删除stream失败", err)) }
	ErrNatsNotFooundStream  = func(f string) error { return gerror.NewCode(gcode.New(200004, "stream未找到:"+f, nil)) }
)
