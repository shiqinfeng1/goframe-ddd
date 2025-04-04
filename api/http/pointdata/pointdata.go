// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package pointdata

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/http/pointdata/v1"
)

type IPointdataV1 interface {
	PubSubBenchmark(ctx context.Context, req *v1.PubSubBenchmarkReq) (res *v1.PubSubBenchmarkRes, err error)
	StreamSend(ctx context.Context, req *v1.StreamSendReq) (res *v1.StreamSendRes, err error)
	GetStreamInfo(ctx context.Context, req *v1.GetStreamInfoReq) (res *v1.GetStreamInfoRes, err error)
	DeleteStream(ctx context.Context, req *v1.DeleteStreamReq) (res *v1.DeleteStreamRes, err error)
}
