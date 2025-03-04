package filemgr

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
)

// FileSendQueue 表示文件发送队列
type FileSendQueue struct {
	maxTasks int
	tasks    []*TransferTask
	running  int
	mu       sync.Mutex
	cond     *sync.Cond
}

// NewFileSendQueue 创建一个新的文件发送队列
func NewfileTransferService(maxTasks int) *FileSendQueue {
	q := &FileSendQueue{
		maxTasks: maxTasks,
		tasks:    make([]*TransferTask, 0),
		running:  0,
	}
	q.cond = sync.NewCond(&q.mu)
	q.start()
	return q
}

// AddTask 向队列中添加一个新的文件发送任务
func (q *FileSendQueue) AddTask(ctx context.Context, id, name string, path []string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	task := &TransferTask{
		ID:     id,
		Path:   path,
		Name:   name,
		Status: StatusWaiting,
	}
	q.tasks = append(q.tasks, task)
	q.cond.Signal()
}

// Start 开始处理队列中的任务
func (q *FileSendQueue) start() {
	go func() {
		for {
			q.mu.Lock()
			for q.running >= q.maxTasks || len(q.tasks) == 0 {
				q.cond.Wait()
			}

			// 找到第一个等待发送的任务
			var nextTask *TransferTask
			for _, task := range q.tasks {
				if task.Status == StatusWaiting {
					nextTask = task
					break
				}
			}

			if nextTask != nil {
				nextTask.Status = StatusSending
				q.running++
				q.mu.Unlock()

				// 模拟文件发送
				// err := q.sendFile(nextTask)
				err := gerror.New("d")
				q.mu.Lock()
				if err != nil {
					nextTask.Status = StatusFailed
				} else {
					nextTask.Status = StatusSuccessful
				}
				q.running--
				q.cond.Signal()
			}
			q.mu.Unlock()
		}
	}()
}

// PauseTask 暂停指定 ID 的任务
func (q *FileSendQueue) PauseTask(ctx context.Context, id string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, task := range q.tasks {
		if task.ID == id && task.Status == StatusSending {
			task.Status = StatusPaused
			break
		}
	}
}

// ResumeTask 恢复指定 ID 的任务
func (q *FileSendQueue) ResumeTask(ctx context.Context, id string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, task := range q.tasks {
		if task.ID == id && task.Status == StatusPaused {
			task.Status = StatusWaiting
			q.cond.Signal()
			break
		}
	}
}

// GetTaskStatus 获取指定 ID 任务的状态
func (q *FileSendQueue) GetTaskStatus(ctx context.Context, id string) string {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, task := range q.tasks {
		if task.ID == id {
			return task.Status
		}
	}
	return ""
}
