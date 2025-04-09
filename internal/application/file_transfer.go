package application

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/rs/xid"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
)

func (app *Application) StartSendFile(ctx context.Context, in *StartSendFileInput) (*StartSendFileOutput, error) {
	taskId := xid.New().String()
	app.fileTransfer.AddTask(ctx, taskId, in.BaseName, in.NodeId, in.Files)
	return nil, nil
}

func (app *Application) PauseSendFile(ctx context.Context, in *PauseSendFileInput) (*PauseSendFileOutput, error) {
	app.fileTransfer.PauseTask(ctx, in.TaskId)
	return nil, nil
}

func (app *Application) CancelSendFile(ctx context.Context, in *CancelSendFileInput) (*CancelSendFileOutput, error) {
	app.fileTransfer.CancelTask(ctx, in.TaskId)
	return nil, nil
}

func (app *Application) ResumeSendFile(ctx context.Context, in *ResumeSendFileInput) (*ResumeSendFileOutput, error) {
	app.fileTransfer.ResumeTask(ctx, in.TaskId)
	return nil, nil
}

func (app *Application) RemoveTask(ctx context.Context, in *RemoveTaskInput) (*RemoveTaskOutput, error) {
	app.fileTransfer.RemoveTask(ctx, in.TaskIds)
	return nil, nil
}

func (app *Application) GetClientIds(ctx context.Context) ([]string, error) {
	nodeIds, err := filemgr.Session().GetNodeList(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, nil
	}
	return nodeIds, nil
}

func (app *Application) GetSendingTaskList(ctx context.Context, in *TaskListInput) (*TaskListOutput, error) {
	running, maxTasks := app.fileTransfer.GetMaxAndRunning(ctx)
	tasks, sfs, err := app.fileTransfer.GetNotCompletedTasks(ctx)
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

func (app *Application) GetCompletedTaskList(ctx context.Context, in *TaskListInput) (*TaskListOutput, error) {
	tasks, sfs, err := app.fileTransfer.GetCompletedTasks(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, nil
	}
	tasklist := &TaskListOutput{}
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
			SendedPercent: fmt.Sprintf("%.2f", 100*sended/sendTotal),
			Speed:         task.Speed,
			Elapsed:       task.Elapsed,
		})
	}

	return tasklist, nil
}
