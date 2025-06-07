// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package ops

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/ops/v1"
)

type IOpsV1 interface {
	ImageList(ctx context.Context, req *v1.ImageListReq) (res *v1.ImageListRes, err error)
	ContainerImage(ctx context.Context, req *v1.ContainerImageReq) (res *v1.ContainerImageRes, err error)
	GetStreamList(ctx context.Context, req *v1.GetStreamListReq) (res *v1.GetStreamListRes, err error)
	GetStreamInfo(ctx context.Context, req *v1.GetStreamInfoReq) (res *v1.GetStreamInfoRes, err error)
	DeleteStream(ctx context.Context, req *v1.DeleteStreamReq) (res *v1.DeleteStreamRes, err error)
	UpgradeApp(ctx context.Context, req *v1.UpgradeAppReq) (res *v1.UpgradeAppRes, err error)
	UpgradeImage(ctx context.Context, req *v1.UpgradeImageReq) (res *v1.UpgradeImageRes, err error)
}
