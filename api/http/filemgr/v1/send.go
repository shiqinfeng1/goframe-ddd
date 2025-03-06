package v1

import "github.com/gogf/gf/v2/frame/g"

type StartSendFileReq struct {
	g.Meta   `path:"/file/startSend" tags:"文件收发" method:"post" summary:"开始发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件/目录绝对路径"`
	NodeId   string   `p:"node_id" dc:"服务端发送时需要传入指定的客户端节点id"`
}
type StartSendFileRes struct {
	g.Meta `status:"200"`
}

type PauseSendFileReq struct {
	g.Meta   `path:"/file/pauseSend" tags:"文件收发" method:"post" summary:"暂停发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件绝对路径"`
}
type PauseSendFileRes struct {
	g.Meta `status:"200"`
}

type CancelSendFileReq struct {
	g.Meta   `path:"/file/cancelSend" tags:"文件收发" method:"post" summary:"取消发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件绝对路径"`
}
type CancelSendFileRes struct {
	g.Meta `status:"200"`
}
type ResumeSendFileReq struct {
	g.Meta   `path:"/file/resume" tags:"文件收发" method:"post" summary:"继续发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件绝对路径"`
}
type ResumeSendFileRes struct {
	g.Meta `status:"200"`
}
