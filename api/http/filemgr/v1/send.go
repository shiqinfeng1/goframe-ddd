package v1

import "github.com/gogf/gf/v2/frame/g"

type StartSendFileReq struct {
	g.Meta   `path:"/file/startSend" tags:"文件收发" method:"post" summary:"开始发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件/目录绝对路径"`
}
type StartSendFileRes struct {
	g.Meta `mime:"application/json"`
}

type PauseSendFileReq struct {
	g.Meta   `path:"/file/pauseSend" tags:"文件收发" method:"post" summary:"暂停发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件绝对路径"`
}
type PauseSendFileRes struct {
	g.Meta `mime:"application/json"`
}

type CancelSendFileReq struct {
	g.Meta   `path:"/file/cancelSend" tags:"文件收发" method:"post" summary:"取消发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件绝对路径"`
}
type CancelSendFileRes struct {
	g.Meta `mime:"application/json"`
}
type ResumeSendFileReq struct {
	g.Meta   `path:"/file/resume" tags:"文件收发" method:"post" summary:"继续发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件绝对路径"`
}
type ResumeSendFileRes struct {
	g.Meta `mime:"application/json"`
}
