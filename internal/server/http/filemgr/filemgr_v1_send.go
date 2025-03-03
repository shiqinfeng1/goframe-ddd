package filemgr

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/filemgr/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/application/command"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) StartSendFile(ctx context.Context, req *v1.StartSendFileReq) (res *v1.StartSendFileRes, err error) {
	for _, v := range req.FilePath {
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
				return nil, errors.ErrInvalidFile(gconv.String(req.FilePath))
			}
			_, err = c.app.Commands.StartSendFile(ctx, &command.StartSendFileInput{
				IsDir:    true,
				BaseName: gfile.Basename(v),
				Files:    files,
			})
		} else { // 发送文件
			var files []string
			if !gfile.IsEmpty(v) {
				files = append(files, v)
			}
			if len(files) == 0 {
				return nil, errors.ErrInvalidFile(gconv.String(req.FilePath))
			}
			_, err = c.app.Commands.StartSendFile(ctx, &command.StartSendFileInput{
				IsDir:    true,
				BaseName: gfile.Basename(v),
				Files:    files,
			})
		}
	}

	return
}

func (c *ControllerV1) PauseSendFile(ctx context.Context, req *v1.PauseSendFileReq) (res *v1.PauseSendFileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (c *ControllerV1) CancelSendFile(ctx context.Context, req *v1.CancelSendFileReq) (res *v1.CancelSendFileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (c *ControllerV1) ResumeSendFile(ctx context.Context, req *v1.ResumeSendFileReq) (res *v1.ResumeSendFileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
