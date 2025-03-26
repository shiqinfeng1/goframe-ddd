package v1

import "github.com/gogf/gf/v2/frame/g"

type PubSubBenchmarkReq struct {
	g.Meta  `path:"/pubsub/benchmark" tags:"消息队列测试" method:"post" summary:"消息发布基准测试"`
	NumPubs int `p:"num_pubs" dc:"并发的发布者数量. 默认值:1"`
	NumSubs int `p:"num_subs" dc:"并发的订阅者数量. 默认值:1"`
	NumMsgs int `p:"num_msgs" dc:"发布的消息个数. 默认值:100000"`
	MsgSize int `p:"msg_size" dc:"每个消息的尺寸. 默认值:128"`
}
type PubSubBenchmarkRes struct {
	g.Meta `status:"200"`
}
