package filemgr

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// FileTransferMgr 表示文件发送队列
type FileTransferMgr struct {
	running  int
	maxTasks int
	tasks    gmap.StrAnyMap
	mutex    sync.Mutex
	notify   chan struct{}
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
	ctx := gctx.New()
	// 从数据库中取出未完成的任务
	tasks, sendfiles, err := q.repo.GetNotCompletedTasks(ctx)
	if err != nil {
		g.Log().Fatalf(ctx, "get not completed task fail:%v", err)
	}

	var paths []string
	for _, task := range tasks {
		for _, sf := range sendfiles[task.TaskID] {
			paths = append(paths, sf.FilePath)
		}
		newtask := NewTransferTask(ctx, task.TaskID, task.TaskName, task.NodeID, paths, NewStatus(uint32(task.Status)), q.stream, q.repo)
		// 缓存到tasks队列
		q.tasks.Set(task.TaskID, newtask)
		// todo 何时触发重新发送未完成的任务?
	}

	q.notify = make(chan struct{}, 16) // 运行多个任务通知缓存到通道里
	q.start(ctx)
	return q
}

func (q *FileTransferMgr) GetMaxAndRunning(ctx context.Context) (int, int) {
	return q.running, q.maxTasks
}

func (q *FileTransferMgr) GetNotCompletedTasks(ctx context.Context) ([]*FileTransferTask, map[string][]*SendFile, error) {
	return q.repo.GetNotCompletedTasks(ctx)
}

func (q *FileTransferMgr) GetCompletedTasks(ctx context.Context) ([]*FileTransferTask, map[string][]*SendFile, error) {
	return q.repo.GetCompletedTasks(ctx)
}

// AddTask 向队列中添加一个新的文件发送任务
func (q *FileTransferMgr) AddTask(ctx context.Context, id, name, nodeId string, paths []string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	task := NewTransferTask(ctx, id, name, nodeId, paths, StatusWaiting, q.stream, q.repo)
	q.tasks.Set(id, task)
	q.notify <- struct{}{}
}

// Start 开始处理队列中的任务
func (q *FileTransferMgr) start(ctx context.Context) {
	go func() {
		for {
			<-q.notify // 等待通知
			q.mutex.Lock()
			if q.running >= q.maxTasks {
				g.Log().Infof(ctx, "%v task is running, waitting...", q.running)
				q.mutex.Unlock()
				continue
			}
			if q.tasks.Size() == 0 {
				g.Log().Infof(ctx, "no task in queue!")
				q.mutex.Unlock()
				continue
			}
			// 找到第一个等待发送的任务
			q.tasks.Iterator(func(k string, v any) bool {
				task := v.(*TransferTask)
				if task.status == StatusWaiting {
					// 更新状态
					task.status = StatusSending
					q.running++
					task.sendFileChan <- func(success bool) {
						q.mutex.Lock()
						if success {
							q.tasks.Remove(k) // 直接移除 不需要更新状态为StatusSuccessful
							g.Log().Debugf(ctx, "task:%v(%v) finished success: %v!", task.taskId, task.taskName, StatusSuccessful)
						} else {
							task.status = StatusFailed
							g.Log().Debugf(ctx, "task:%v(%v) finished fail: %v!", task.taskId, task.taskName, task.status)
						}
						q.running--
						q.mutex.Unlock()
						q.notify <- struct{}{}
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
func (q *FileTransferMgr) PauseTask(ctx context.Context, id string) {
	var (
		task      *TransferTask
		oldStatus Status
	)

	q.mutex.Lock()
	q.tasks.Iterator(func(k string, v any) bool {
		task = v.(*TransferTask)
		if task.taskId == id &&
			(task.status == StatusSending || task.status == StatusWaiting) {
			oldStatus = task.status
			task.pauseChan <- func() {
				q.mutex.Lock()
				task.status = StatusPaused
				q.running--
				q.mutex.Unlock()
			}
			return false
		}
		return true
	})
	q.mutex.Unlock()
	g.Log().Debugf(ctx, "task:%v(%v) change status  %v -> %v !", task.taskId, task.taskName, oldStatus, task.status)
	q.notify <- struct{}{}
}

// CancelTask 取消指定 ID 的任务
func (q *FileTransferMgr) CancelTask(ctx context.Context, id string) {
	var (
		needRemove bool
		task       *TransferTask
		oldStatus  Status
	)

	q.mutex.Lock()
	q.tasks.Iterator(func(k string, v any) bool {
		task = v.(*TransferTask)
		if task.taskId == id &&
			(task.status == StatusSending || task.status == StatusWaiting || task.status == StatusPaused) {

			oldStatus = task.status
			task.status = StatusPaused
			task.cancelChan <- func() {
				q.mutex.Lock()
				task.status = StatusCancel
				q.running--
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
	q.mutex.Unlock()
	if needRemove {
		g.Log().Debugf(ctx, "task:%v(%v) change status  %v -> %v !", task.taskId, task.taskName, oldStatus, task.status)
		q.notify <- struct{}{}
	}
}

// ResumeTask 恢复指定 ID 的任务
func (q *FileTransferMgr) ResumeTask(ctx context.Context, id string) {
	var (
		found bool
		task  *TransferTask
	)
	q.mutex.Lock()
	q.tasks.Iterator(func(k string, v any) bool {
		task = v.(*TransferTask)
		if task.taskId == id && (task.status == StatusPaused || task.status == StatusFailed) {
			task.status = StatusWaiting
			found = true
			return false
		}
		return true
	})
	q.mutex.Unlock()
	if found {
		g.Log().Debugf(ctx, "task:%v(%v) change status %v -> %v !", task.taskId, task.taskName, StatusPaused, task.status)
		q.notify <- struct{}{}
	}
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

func (q *FileTransferMgr) RemoveTask(ctx context.Context, taskIds []string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for _, taskId := range taskIds {
		if q.tasks.Contains(taskId) {
			old := q.tasks.Remove(taskId)
			task := old.(*TransferTask)
			task.Exit()
		}
	}
	if err := q.repo.RemoveTasks(ctx, taskIds); err != nil {
		g.Log().Errorf(ctx, "remove task:%v fail", taskIds)
	}
}
