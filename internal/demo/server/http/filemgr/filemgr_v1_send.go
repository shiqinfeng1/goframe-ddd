package filemgr

import (
	"context"
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/demo/http/filemgr/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/demo/application/dto"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) StartSendFile(ctx context.Context, req *v1.StartSendFileReq) (res *v1.StartSendFileRes, err error) {
	res = &v1.StartSendFileRes{}
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	if g.Cfg().MustGet(ctx, "sessionmgr.isCloud").Bool() {
		if req.NodeId == "" {
			return res, errors.ErrInvalidNodeId(lang)
		}
	} else {
		req.NodeId = ""
	}
	for _, v := range req.FilePath {
		if !filepath.IsAbs(v) {
			return res, gerror.Wrap(errors.ErrNotAbsFilePath(lang), v)
		}
		// 发送文件夹
		if gfile.IsDir(v) {
			var files []string
			gfile.ScanDirFileFunc(v, "*", true, func(p string) string {
				if !gfile.IsEmpty(p) {
					files = append(files, p)
				}
				return ""
			})
			if len(files) == 0 {
				return res, gerror.Wrap(errors.ErrEmptyDir(lang), v)
			}
			c.app.StartSendFile(ctx, &dto.StartSendFileInput{
				BaseName: gfile.Basename(v),
				Files:    files,
				NodeId:   req.NodeId,
			})
		} else { // 发送文件
			var files []string
			if !gfile.IsEmpty(v) {
				files = append(files, v)
			}
			if len(files) == 0 {
				return res, gerror.Wrap(errors.ErrInvalidFiles(lang), gconv.String(req.FilePath))
			}
			c.app.StartSendFile(ctx, &dto.StartSendFileInput{
				BaseName: gfile.Basename(v),
				Files:    files,
				NodeId:   req.NodeId,
			})
		}
	}

	return res, nil
}

// func (c *ControllerV1) PauseSendFile(ctx context.Context, req *v1.PauseSendFileReq) (res *v1.PauseSendFileRes, err error) {
// 	res = &v1.PauseSendFileRes{}

// 	_, err = c.app.PauseSendFile(ctx, &application.PauseSendFileInput{
// 		TaskId: req.TaskId,
// 	})
// 	return
// }

// func (c *ControllerV1) CancelSendFile(ctx context.Context, req *v1.CancelSendFileReq) (res *v1.CancelSendFileRes, err error) {
// 	res = &v1.CancelSendFileRes{}

// 	_, err = c.app.CancelSendFile(ctx, &application.CancelSendFileInput{
// 		TaskId: req.TaskId,
// 	})
// 	return
// }

// func (c *ControllerV1) ResumeSendFile(ctx context.Context, req *v1.ResumeSendFileReq) (res *v1.ResumeSendFileRes, err error) {
// 	res = &v1.ResumeSendFileRes{}

// 	_, err = c.app.ResumeSendFile(ctx, &application.ResumeSendFileInput{
// 		TaskId: req.TaskId,
// 	})
// 	return
// }

// func (c *ControllerV1) RemoveTask(ctx context.Context, req *v1.RemoveTaskReq) (res *v1.RemoveTaskRes, err error) {
// 	res = &v1.RemoveTaskRes{}

// 	_, err = c.app.RemoveTask(ctx, &application.RemoveTaskInput{
// 		TaskIds: req.TaskIds,
// 	})
// 	return
// }

func (c *ControllerV1) SendingTaskList(ctx context.Context, req *v1.SendingTaskListReq) (res *v1.SendingTaskListRes, err error) {

	out, _ := c.app.GetSendingTaskList(ctx, &dto.TaskListInput{})
	res = &v1.SendingTaskListRes{
		Running:  out.Running,
		MaxTasks: out.MaxTasks,
		Tasks:    out.Tasks,
	}
	return
}

func (c *ControllerV1) CompletedTaskList(ctx context.Context, req *v1.CompletedTaskListReq) (res *v1.CompletedTaskListRes, err error) {

	out, _ := c.app.GetCompletedTaskList(ctx, &dto.TaskListInput{})
	res = &v1.CompletedTaskListRes{
		Tasks: out.Tasks,
	}
	return
}
