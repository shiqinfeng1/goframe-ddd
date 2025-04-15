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
	SendingTaskList(ctx context.Context, req *v1.SendingTaskListReq) (res *v1.SendingTaskListRes, err error)
	CompletedTaskList(ctx context.Context, req *v1.CompletedTaskListReq) (res *v1.CompletedTaskListRes, err error)
}
