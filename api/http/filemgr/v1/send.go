package v1

import "github.com/gogf/gf/v2/frame/g"

type SendFileReq struct {
	g.Meta   `path:"/file/send" tags:"文件收发" method:"post" summary:"发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件绝对路径"`
}
type SendFileRes struct {
	g.Meta `mime:"application/json"`
}
