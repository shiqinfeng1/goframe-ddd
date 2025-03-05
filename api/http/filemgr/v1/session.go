package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
)

type NodeListReq struct {
	g.Meta `path:"/node/list" tags:"会话管理" method:"post" summary:"查询已连接节点列表"`
}
type NodeListRes struct {
	g.Meta  `mime:"application/json"`
	NodeIds []string `json:"node_ids" dc:"已连接节点客户端id列表"`
}

func (r NodeListRes) EnhanceResponseStatus() (resList map[int]goai.EnhancedStatusType) {
	return map[int]goai.EnhancedStatusType{
		200: {
			Examples: ghttp.DefaultHandlerResponse{
				Code:    0,
				Message: "",
				Data:    r,
			},
		},
	}
}
