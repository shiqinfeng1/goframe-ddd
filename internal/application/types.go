package application

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

type DeleteStreamInput struct {
	Name string
}
type PubSubStreamInfoInput struct {
}

type PubSubStreamInfoOutput struct {
	StreamInfo    *jetstream.StreamInfo
	ConsumerInfos []*jetstream.ConsumerInfo
}

type PubSubBenchmarkInput struct {
	NumPubs  int
	NumSubs  int
	NumMsgs  int
	MsgSize  int
	Subjects []string
	Typ      string
}

type File struct {
	Files []string
}
type StartSendFileInput struct {
	NodeId   string
	BaseName string
	Files    []string // 如果是目录，那么是该目录下的所有文件列表
}
type StartSendFileOutput struct{}

type ClientIdsOutput struct {
	Ids []string
}

type PauseSendFileInput struct {
	TaskId string
}
type PauseSendFileOutput struct{}

type CancelSendFileInput struct {
	TaskId string
}
type CancelSendFileOutput struct{}

type ResumeSendFileInput struct {
	TaskId string
}
type ResumeSendFileOutput struct{}

type RemoveTaskInput struct {
	TaskIds []string
}
type RemoveTaskOutput struct{}

type Task struct {
	TaskName      string   `json:"task_name" dc:"任务名称"`
	TaskId        string   `json:"task_id" dc:"任务id"`
	NodeId        string   `json:"node_id" dc:"节点id"`
	Paths         []string `json:"paths" dc:"发送文件路径"`
	Status        int      `json:"status" dc:"任务状态 1:等待发送 2:正在发送 3:已暂停 4:已取消 5:发送失败 6:发送成功"`
	SendedPercent string   `json:"sended_percent" dc:"发送百分比"`
	Speed         string   `json:"speed" dc:"速率"`
	Elapsed       string   `json:"elapsed" dc:"耗时"`
}

type TaskListInput struct{}

type TaskListOutput struct {
	Running  int `json:"runnings" dc:"正在运行的发送任务数量"`
	MaxTasks int `json:"max_tasks" dc:"同时运行的最大发送任务数量"`
	Tasks    []Task
}
