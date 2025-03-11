package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/application/query"
)

type StartSendFileReq struct {
	g.Meta   `path:"/file/startSend" tags:"文件收发" method:"post" summary:"开始发送文件"`
	FilePath []string `p:"file_path" v:"required" dc:"文件/目录绝对路径"`
	NodeId   string   `p:"node_id" dc:"服务端发送时需要传入指定的客户端节点id"`
}
type StartSendFileRes struct {
	g.Meta `status:"200"`
}

type PauseSendFileReq struct {
	g.Meta `path:"/file/pauseSend" tags:"文件收发" method:"post" summary:"暂停发送文件"`
	TaskId string `p:"task_id" v:"required" dc:"任务id"`
}
type PauseSendFileRes struct {
	g.Meta `status:"200"`
}

type CancelSendFileReq struct {
	g.Meta `path:"/file/cancelSend" tags:"文件收发" method:"post" summary:"取消发送文件"`
	TaskId string `p:"task_id" v:"required" dc:"任务id"`
}
type CancelSendFileRes struct {
	g.Meta `status:"200"`
}
type ResumeSendFileReq struct {
	g.Meta `path:"/file/resumeSend" tags:"文件收发" method:"post" summary:"继续发送文件"`
	TaskId string `p:"task_id" v:"required" dc:"任务id"`
}
type ResumeSendFileRes struct {
	g.Meta `status:"200"`
}

type SendingTaskListReq struct {
	g.Meta `path:"/task/sendingList" tags:"文件收发" method:"post" summary:"查询未完成的任务列表"`
}

type SendingTaskListRes struct {
	g.Meta   `status:"200"`
	Running  int          `json:"runnings" dc:"正在运行的发送任务数量"`
	MaxTasks int          `json:"max_tasks" dc:"同时运行的最大发送任务数量"`
	Tasks    []query.Task `json:"tasks" dc:"已连接节点客户端id列表"`
}

type CompletedTaskListReq struct {
	g.Meta `path:"/task/completedList" tags:"文件收发" method:"post" summary:"查询已完成的任务列表"`
}

type CompletedTaskListRes struct {
	g.Meta `status:"200"`
	Tasks  []query.Task `json:"tasks" dc:"已连接节点客户端id列表"`
}

type RemoveTaskReq struct {
	g.Meta  `path:"/task/remove" tags:"文件收发" method:"post" summary:"删除任务"`
	TaskIds []string `p:"task_ids" v:"required" dc:"任务id列表"`
}

type RemoveTaskRes struct {
	g.Meta `status:"200"`
}
