package filemgr

// 任务状态枚举
type Status struct {
	val  int
	desc string
}

var (
	StatusUndefined  = Status{val: 0, desc: "未定义"}
	StatusWaiting    = Status{val: 1, desc: "等待发送"}
	StatusSending    = Status{val: 2, desc: "正在发送"}
	StatusPaused     = Status{val: 3, desc: "已暂停"}
	StatusCancel     = Status{val: 4, desc: "已取消"}
	StatusFailed     = Status{val: 5, desc: "发送失败"}
	StatusSuccessful = Status{val: 6, desc: "发送成功"}
)

type (
	postSendFunc func(bool)
	postFunc     func()
)

// FileSendTask 表示文件发送任务
type TransferTask struct {
	name       string
	id         string
	paths      []string
	status     Status
	sendChan   chan postSendFunc
	pauseChan  chan postFunc
	cancelChan chan postFunc
}
