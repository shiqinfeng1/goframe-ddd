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
	ErrNotAbsFilePath = func(f string) error { return gerror.NewCode(gcode.New(100001, "请输入绝对路径:"+f, nil)) }
	ErrEmptyDir       = func(f string) error { return gerror.NewCode(gcode.New(100002, "文件夹内无有效文件:"+f, nil)) }
	ErrInvalidFiles   = func(f string) error { return gerror.NewCode(gcode.New(100003, "无有效文件:"+f, nil)) }
	ErrInvalidNodeId  = gerror.NewCode(gcode.New(100004, "无效的客户端节点ID", nil))
	ErrFileMgrDisable = gerror.NewCode(gcode.New(100005, "未使能文件收发功能", nil))

	ErrOver4GSize = gerror.NewCode(gcode.New(100101, "文件尺寸不能大于4G", nil))
)
