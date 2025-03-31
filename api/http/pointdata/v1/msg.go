package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
)

type PubSubBenchmarkReq struct {
	g.Meta   `path:"/pubsub/benchmark" tags:"消息队列管理" method:"post" summary:"消息发布订阅基准测试"`
	NumPubs  int      `p:"num_pubs" dc:"并发的发布者数量. 默认值:1"`
	NumSubs  int      `p:"num_subs" dc:"并发的订阅者数量. 默认值:1"`
	NumMsgs  int      `p:"num_msgs" dc:"发布的消息个数. 默认值:100000"`
	MsgSize  int      `p:"msg_size" dc:"每个消息的尺寸. 默认值:128"`
	Subjects []string `p:"subjects" dc:"主题. 可以配置多个. 默认值: benchmark-test"`
	Typ      string   `p:"typ" v:"required|in:pubsub,jetstream" dc:"测试类型. 默认值: pubsub"`
}
type PubSubBenchmarkRes struct {
	g.Meta `status:"200"`
}

type GetStreamInfoReq struct {
	g.Meta `path:"/pubsub/stream/getInfo" tags:"消息队列管理" method:"post" summary:"查询消息流信息和状态"`
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
