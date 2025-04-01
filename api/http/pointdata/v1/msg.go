package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
)

type PubSubBenchmarkReq struct {
	g.Meta  `path:"/pubsub/benchmark" tags:"消息队列管理" method:"post" summary:"消息发布订阅基准测试"`
	MsgSize int `p:"msg_size" dc:"每个消息的尺寸. 默认值:128"`
}
type PubSubBenchmarkRes struct {
	g.Meta `status:"200"`
}

type GetStreamInfoReq struct {
	g.Meta     `path:"/pubsub/stream/getInfo" tags:"消息队列管理" method:"post" summary:"查询消息流信息和状态"`
	StreamName string `p:"stream_name" v:"required" dc:"消息流名称"`
}
type GetStreamInfoRes struct {
	g.Meta        `status:"200"`
	StreamInfo    *application.StreamInfo     `json:"stream_info" dc:"流信息"`
	ConsumerInfos []*application.ConsumerInfo `json:"consumer_infos" dc:"流对应的消费者信息"`
}

type DeleteStreamReq struct {
	g.Meta     `path:"/pubsub/stream/delete" tags:"消息队列管理" method:"post" summary:"删除消息流"`
	StreamName string `p:"stream_name" dc:"消息流名称"`
}
type DeleteStreamRes struct {
	g.Meta `status:"200"`
}
