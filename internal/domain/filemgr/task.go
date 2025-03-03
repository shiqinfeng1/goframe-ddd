package filemgr

// 任务状态枚举
const (
	StatusWaiting    = "等待发送"
	StatusSending    = "正在发送"
	StatusPaused     = "已暂停"
	StatusFailed     = "发送失败"
	StatusSuccessful = "发送成功"
)

// FileSendTask 表示文件发送任务
type TransferTask struct {
	Name   string
	ID     string
	Path   []string
	Status string
}
