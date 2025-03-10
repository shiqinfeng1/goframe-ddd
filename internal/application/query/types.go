package query

type Task struct {
	TaskName      string   `json:"task_name" dc:"任务名称"`
	TaskId        string   `json:"task_id" dc:"任务id"`
	NodeId        string   `json:"node_id" dc:"节点id"`
	Paths         []string `json:"paths" dc:"正在运行的发送任务数量"`
	Status        int      `json:"status" dc:"任务状态 1:等待发送 2:正在发送 3:已暂停 4:已取消 5:发送失败 6:发送成功"`
	SendedPercent string   `json:"sended_percent" dc:"正在运行的发送任务数量"`
}

type TaskListInput struct{}

type TaskListOutput struct {
	Running  int `json:"runnings" dc:"正在运行的发送任务数量"`
	MaxTasks int `json:"max_tasks" dc:"同时运行的最大发送任务数量"`
	Tasks    []Task
}
