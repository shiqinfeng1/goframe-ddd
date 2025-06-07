package errors

import "github.com/gogf/gf/v2/errors/gcode"

// code格式： xxyyzz    xx：服务  yy:模块  zz：错误编号
// 10：文件传输服务
// --00：发送模块
// --01: 文件读写模块
// 20：消息队列
// --00：nats管理
var (
	ErrNotAbsFilePath = codeFunc(gcode.New(CodeNotAbsFilePath, "", nil), "file.notAbsFilePath")
	ErrEmptyDir       = codeFunc(gcode.New(CodeEmptyDir, "", nil), "file.emptyDir")
	ErrInvalidFiles   = codeFunc(gcode.New(CodeInvalidFiles, "", nil), "file.invalidFiles")
	ErrInvalidNodeId  = codeFunc(gcode.New(CodeInvalidNodeId, "", nil), "file.invalidNodeId")
	ErrOver4GSize     = codeFunc(gcode.New(CodeOver4GSize, "", nil), "file.over4GSize")
)
