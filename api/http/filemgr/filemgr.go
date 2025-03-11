// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package filemgr

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/http/filemgr/v1"
)

type IFilemgrV1 interface {
	StartSendFile(ctx context.Context, req *v1.StartSendFileReq) (res *v1.StartSendFileRes, err error)
	PauseSendFile(ctx context.Context, req *v1.PauseSendFileReq) (res *v1.PauseSendFileRes, err error)
	CancelSendFile(ctx context.Context, req *v1.CancelSendFileReq) (res *v1.CancelSendFileRes, err error)
	ResumeSendFile(ctx context.Context, req *v1.ResumeSendFileReq) (res *v1.ResumeSendFileRes, err error)
	SendingTaskList(ctx context.Context, req *v1.SendingTaskListReq) (res *v1.SendingTaskListRes, err error)
	CompletedTaskList(ctx context.Context, req *v1.CompletedTaskListReq) (res *v1.CompletedTaskListRes, err error)
	RemoveTask(ctx context.Context, req *v1.RemoveTaskReq) (res *v1.RemoveTaskRes, err error)
	NodeList(ctx context.Context, req *v1.NodeListReq) (res *v1.NodeListRes, err error)
}
