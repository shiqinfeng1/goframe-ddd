// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package pointdata

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/pointdata/v1"
)

type IPointdataV1 interface {
	GetStreamList(ctx context.Context, req *v1.GetStreamListReq) (res *v1.GetStreamListRes, err error)
	GetStreamInfo(ctx context.Context, req *v1.GetStreamInfoReq) (res *v1.GetStreamInfoRes, err error)
	DeleteStream(ctx context.Context, req *v1.DeleteStreamReq) (res *v1.DeleteStreamRes, err error)
}
