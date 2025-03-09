package filemgr

import (
	"context"
	"path/filepath"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/filemgr/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/application/command"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) StartSendFile(ctx context.Context, req *v1.StartSendFileReq) (res *v1.StartSendFileRes, err error) {
	res = &v1.StartSendFileRes{}
	if g.Cfg().MustGet(ctx, "filemgr.isCloud").Bool() {
		if req.NodeId == "" {
			return res, errors.ErrInvalidNodeId
		}
	} else {
		req.NodeId = ""
	}
	for _, v := range req.FilePath {
		if !filepath.IsAbs(v) {
			return res, errors.ErrNotAbsFilePath(v)
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
				return res, errors.ErrEmptyDir(v)
			}
			c.app.Commands.StartSendFile(ctx, &command.StartSendFileInput{
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
				return res, errors.ErrInvalidFiles(gconv.String(req.FilePath))
			}
			c.app.Commands.StartSendFile(ctx, &command.StartSendFileInput{
				BaseName: gfile.Basename(v),
				Files:    files,
				NodeId:   req.NodeId,
			})
		}
	}

	return res, nil
}

func (c *ControllerV1) PauseSendFile(ctx context.Context, req *v1.PauseSendFileReq) (res *v1.PauseSendFileRes, err error) {
	res = &v1.PauseSendFileRes{}
	_, err = c.app.Commands.PauseSendFile(ctx, &command.PauseSendFileInput{
		TaskId: req.TaskId,
	})
	return
}

func (c *ControllerV1) CancelSendFile(ctx context.Context, req *v1.CancelSendFileReq) (res *v1.CancelSendFileRes, err error) {
	res = &v1.CancelSendFileRes{}
	_, err = c.app.Commands.CancelSendFile(ctx, &command.CancelSendFileInput{
		TaskId: req.TaskId,
	})
	return
}

func (c *ControllerV1) ResumeSendFile(ctx context.Context, req *v1.ResumeSendFileReq) (res *v1.ResumeSendFileRes, err error) {
	res = &v1.ResumeSendFileRes{}
	_, err = c.app.Commands.ResumeSendFile(ctx, &command.ResumeSendFileInput{
		TaskId: req.TaskId,
	})
	return
}
