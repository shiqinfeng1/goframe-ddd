package v1

import "github.com/gogf/gf/v2/frame/g"

type SessionListReq struct {
	g.Meta `path:"/session/list" tags:"会话管理" method:"post" summary:"服务端查询已连接客户端会话列表"`
}
type SessionListRes struct {
	g.Meta    `mime:"application/json"`
	ClientIds []string `json:"client_ids" dc:"已连接客户端id列表"`
}
