package query

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
)

type Handler struct {
	fileTransfer FileTransferService
}

func NewHandler(
	fileTransfer FileTransferService,
) *Handler {
	return &Handler{
		fileTransfer: fileTransfer,
	}
}

func (h *Handler) GetClientIds(ctx context.Context) ([]string, error) {
	nodeIds, err := filemgr.Session().GetNodeList(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, nil
	}
	return nodeIds, nil
}

func (h *Handler) GetTaskList(ctx context.Context, in *TaskListInput) (*TaskListOutput, error) {
	running, maxTasks := h.fileTransfer.GetMaxAndRunning(ctx)
	tasks, sfs, err := h.fileTransfer.GetNotCompletedTasks(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, nil
	}
	tasklist := &TaskListOutput{
		Running:  running,
		MaxTasks: maxTasks,
	}
	for _, task := range tasks {
		var (
			paths             []string
			sendTotal, sended float32
		)
		for _, sf := range sfs[task.TaskID] {
			paths = append(paths, sf.FilePath)
			sendTotal += float32(sf.ChunkNumTotal)
			sended += float32(sf.ChunkNumSended)
		}
		tasklist.Tasks = append(tasklist.Tasks, Task{
			TaskName:      task.TaskName,
			TaskId:        task.TaskID,
			NodeId:        task.NodeID,
			Paths:         paths,
			Status:        task.Status, // 任务状态 1:等待发送 2:正在发送 3:已暂停 4:已取消 5:发送失败 6:发送成功
			SendedPercent: fmt.Sprintf("%.2f", sended/sendTotal),
		})
	}

	return tasklist, nil
}
