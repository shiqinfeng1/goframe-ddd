// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package filemgr

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/http/filemgr/v1"
)

type IFilemgrV1 interface {
	SendFile(ctx context.Context, req *v1.SendFileReq) (res *v1.SendFileRes, err error)
}
