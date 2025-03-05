package filemgr

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/container/gmap"
)

// FileSendQueue 表示文件发送队列
type FileSendQueue struct {
	running  int
	maxTasks int
	tasks    gmap.StrAnyMap
	mutex    sync.Mutex
	cond     *sync.Cond
}

// NewFileSendQueue 创建一个新的文件发送队列
func NewFileTransferService(maxTasks int) *FileSendQueue {
	q := &FileSendQueue{
		maxTasks: maxTasks,
		tasks:    *gmap.NewStrAnyMap(true),
		running:  0,
	}
	q.cond = sync.NewCond(&q.mutex)
	q.start()
	return q
}

// AddTask 向队列中添加一个新的文件发送任务
func (q *FileSendQueue) AddTask(ctx context.Context, id, name string, paths []string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	task := &TransferTask{
		id:         id,
		paths:      paths,
		name:       name,
		status:     StatusWaiting,
		sendChan:   make(chan postSendFunc, 4),
		pauseChan:  make(chan postFunc, 4),
		cancelChan: make(chan postFunc, 4),
	}
	q.tasks.Set(id, task)
}

// Start 开始处理队列中的任务
func (q *FileSendQueue) start() {
	go func() {
		for {
			q.mutex.Lock()
			for q.running >= q.maxTasks || q.tasks.Size() == 0 {
				q.cond.Wait() // 等待被通知，同时unlock q.mutex， 等到通知后，会重新lock
			}
			// 找到第一个等待发送的任务
			q.tasks.Iterator(func(k string, v any) bool {
				task := v.(*TransferTask)
				if task.status == StatusWaiting {
					// 更新状态
					task.status = StatusSending
					q.running++
					task.sendChan <- func(success bool) {
						q.mutex.Lock()
						defer q.mutex.Unlock()
						if success {
							task.status = StatusSuccessful
						} else {
							task.status = StatusFailed
						}
						q.running--
						q.cond.Signal()
					}
					return false
				}
				return true
			})
			q.mutex.Unlock()
		}
	}()
}

// PauseTask 暂停指定 ID 的任务
func (q *FileSendQueue) PauseTask(ctx context.Context, id string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.tasks.Iterator(func(k string, v any) bool {
		task := v.(*TransferTask)
		if task.id == id &&
			(task.status == StatusSending || task.status == StatusWaiting) {
			task.pauseChan <- func() {
				q.mutex.Lock()
				task.status = StatusPaused
				q.running--
				q.cond.Signal()
				q.mutex.Unlock()
			}
			return false
		}
		return true
	})
}

// CancelTask 取消指定 ID 的任务
func (q *FileSendQueue) CancelTask(ctx context.Context, id string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.tasks.Iterator(func(k string, v any) bool {
		task := v.(*TransferTask)
		if task.id == id &&
			(task.status == StatusSending || task.status == StatusWaiting || task.status == StatusPaused) {
			task.status = StatusPaused
			task.cancelChan <- func() {
				q.mutex.Lock()
				task.status = StatusCancel
				q.running--
				q.cond.Signal()
				q.mutex.Unlock()
			}
			return false
		}
		return true
	})
}

// ResumeTask 恢复指定 ID 的任务
func (q *FileSendQueue) ResumeTask(ctx context.Context, id string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.tasks.Iterator(func(k string, v any) bool {
		task := v.(*TransferTask)
		if task.id == id && task.status == StatusPaused {
			task.status = StatusWaiting
			return false
		}
		return true
	})
}

// GetTaskStatus 获取指定 ID 任务的状态
func (q *FileSendQueue) GetTaskStatus(ctx context.Context, id string) Status {
	val, found := q.tasks.Search(id)
	if found {
		task := val.(*TransferTask)
		return task.status
	}
	return StatusUndefined
}

func (q *FileSendQueue) GetTaskList(ctx context.Context) []*TransferTask {
	vals := q.tasks.Values()
	out := make([]*TransferTask, 0)
	for _, v := range vals {
		out = append(out, v.(*TransferTask))
	}
	return out
}
