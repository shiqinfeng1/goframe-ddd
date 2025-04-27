package dto

import (
	"github.com/nats-io/nats.go/jetstream"
)

type ConsumerInfo struct {
	Name           string   `json:"name,omitempty" dc:"消费者名称"`
	Durable        string   `json:"durable_name,omitempty"`
	Description    string   `json:"description,omitempty"`
	DeliverPolicy  string   `json:"deliver_policy"`
	AckPolicy      string   `json:"ack_policy"`
	FilterSubject  string   `json:"filter_subject,omitempty"`
	FilterSubjects []string `json:"filter_subjects,omitempty"`
	NumAckPending  int      `json:"num_ack_pending" dc:"已投递但未确认的消息数量"`
	NumRedelivered int      `json:"num_redelivered" dc:"重新投递但未确认的消息数量"`
	NumWaiting     int      `json:"num_waiting" dc:"在拉取模式下，等待拉取的消费者数量"`
	NumPending     uint64   `json:"num_pending" dc:"未投递的消息数量"`
}
type StreamInfo struct {
	Subjects  []string              `json:"subjects,omitempty" dc:"流的主题列表"`
	Retention string                `json:"retention" dc:"保留策略"`
	State     jetstream.StreamState `json:"state"  dc:"流状态信息"`
}

type DeleteStreamIn struct {
	Name string
}
type JetStreamInfoIn struct {
	Name string
}

type JetStreamInfoOut struct {
	StreamInfo    *jetstream.StreamInfo
	ConsumerInfos []*jetstream.ConsumerInfo
}

type PubSubBenchmarkIn struct {
	Subjects     []string
	JsSubjects   []string
	StreamName   string
	ConsumerName string
}
