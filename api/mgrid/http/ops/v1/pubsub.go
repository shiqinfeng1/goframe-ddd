package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
)

type GetStreamListReq struct {
	g.Meta `path:"/pubsub/stream/list" tags:"消息队列管理" method:"post" summary:"查询消息流列表"`
}
type GetStreamListRes struct {
	g.Meta  `status:"200"`
	Streams []*dto.StreamInfo `json:"streams" dc:"流信息"`
}

type GetStreamInfoReq struct {
	g.Meta     `path:"/pubsub/stream/getInfo" tags:"消息队列管理" method:"post" summary:"查询消息流信息和状态"`
	StreamName string `p:"stream_name" v:"required#未指定消息流名称" dc:"消息流名称"`
}
type GetStreamInfoRes struct {
	g.Meta        `status:"200"`
	StreamInfo    *dto.StreamInfo     `json:"stream_info" dc:"流信息"`
	ConsumerInfos []*dto.ConsumerInfo `json:"consumer_infos" dc:"流对应的消费者信息"`
}

type DeleteStreamReq struct {
	g.Meta     `path:"/pubsub/stream/delete" tags:"消息队列管理" method:"post" summary:"删除消息流"`
	StreamName string `p:"stream_name" dc:"消息流名称"`
}
type DeleteStreamRes struct {
	g.Meta `status:"200"`
}
