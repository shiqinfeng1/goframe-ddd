package filemgr

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/shiqinfeng1/goframe-ddd/api/filexfer/http/filemgr/v1"
)

func (c *ControllerV1) StartSendFile(ctx context.Context, req *v1.StartSendFileReq) (res *v1.StartSendFileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) SendingTaskList(ctx context.Context, req *v1.SendingTaskListReq) (res *v1.SendingTaskListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) CompletedTaskList(ctx context.Context, req *v1.CompletedTaskListReq) (res *v1.CompletedTaskListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
