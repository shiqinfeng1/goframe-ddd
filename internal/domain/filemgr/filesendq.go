package filemgr

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/container/gmap"
)

// FileTransferMgr 表示文件发送队列
type FileTransferMgr struct {
	running  int
	maxTasks int
	tasks    gmap.StrAnyMap
	mutex    sync.Mutex
	cond     *sync.Cond
	stream   StreamIntf
	repo     Repository
}

// NewFileSendQueue 创建一个新的文件发送队列
func NewFileTransferService(maxTasks int, stm StreamIntf, repo Repository) *FileTransferMgr {
	q := &FileTransferMgr{
		maxTasks: maxTasks,
		tasks:    *gmap.NewStrAnyMap(true),
		running:  0,
		stream:   stm,
		repo:     repo,
	}

	q.cond = sync.NewCond(&q.mutex)
	q.start()
	return q
}

// AddTask 向队列中添加一个新的文件发送任务
func (q *FileTransferMgr) AddTask(ctx context.Context, id, name, nodeId string, paths []string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	task := NewTransferTask(ctx, id, name, nodeId, paths, q.stream, q.repo)
	q.tasks.Set(id, task)
	q.cond.Signal()
}

// Start 开始处理队列中的任务
func (q *FileTransferMgr) start() {
	go func() {
		for {
			q.mutex.Lock()
			for q.running >= q.maxTasks || q.tasks.Size() == 0 {
				q.cond.Wait() // 等待被通知，同时unlock q.mutex， 等到通知后，会重新lock
			}
			// 找到第一个等待发送的任务
			var needRemove string
			q.tasks.Iterator(func(k string, v any) bool {
				task := v.(*TransferTask)
				if task.status == StatusWaiting {
					// 更新状态
					task.status = StatusSending
					q.running++
					task.sendFileChan <- func(success bool) {
						q.mutex.Lock()
						defer q.mutex.Unlock()
						if success {
							task.status = StatusSuccessful
							needRemove = k
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
			if needRemove != "" {
				q.tasks.Remove(needRemove)
			}
			q.mutex.Unlock()
		}
	}()
}

// PauseTask 暂停指定 ID 的任务
func (q *FileTransferMgr) PauseTask(ctx context.Context, id string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.tasks.Iterator(func(k string, v any) bool {
		task := v.(*TransferTask)
		if task.taskId == id &&
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
func (q *FileTransferMgr) CancelTask(ctx context.Context, id string) {
	var needRemove bool
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.tasks.Iterator(func(k string, v any) bool {
		task := v.(*TransferTask)
		if task.taskId == id &&
			(task.status == StatusSending || task.status == StatusWaiting || task.status == StatusPaused) {
			task.status = StatusPaused
			task.cancelChan <- func() {
				q.mutex.Lock()
				task.status = StatusCancel
				q.running--
				q.cond.Signal()
				q.mutex.Unlock()
			}
			needRemove = true
			return false
		}
		return true
	})
	if needRemove {
		q.tasks.Remove(id)
	}
}

// ResumeTask 恢复指定 ID 的任务
func (q *FileTransferMgr) ResumeTask(ctx context.Context, id string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.tasks.Iterator(func(k string, v any) bool {
		task := v.(*TransferTask)
		if task.taskId == id && task.status == StatusPaused {
			task.status = StatusWaiting
			q.cond.Signal()
			return false
		}
		return true
	})
}

// GetTaskStatus 获取指定 ID 任务的状态
func (q *FileTransferMgr) GetTaskStatus(ctx context.Context, id string) Status {
	val, found := q.tasks.Search(id)
	if found {
		task := val.(*TransferTask)
		return task.status
	}
	return StatusUndefined
}

func (q *FileTransferMgr) GetTaskList(ctx context.Context) []*TransferTask {
	vals := q.tasks.Values()
	out := make([]*TransferTask, 0)
	for _, v := range vals {
		out = append(out, v.(*TransferTask))
	}
	return out
}
