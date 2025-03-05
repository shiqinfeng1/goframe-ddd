package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type NodeListReq struct {
	g.Meta `path:"/node/list" tags:"会话管理" method:"post" summary:"查询已连接节点列表"`
}
type NodeListRes struct {
	g.Meta  `status:"200"`
	NodeIds []string `json:"node_ids" dc:"已连接节点客户端id列表"`
}
