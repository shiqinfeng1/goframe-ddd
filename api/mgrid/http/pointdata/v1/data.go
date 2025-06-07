package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type GetPointDataReq struct {
	g.Meta `path:"/pointdata/get" tags:"点位数据管理" method:"post" summary:"点位数据查询"`
}
type GetPointDataRes struct {
	g.Meta `status:"200"`
}
