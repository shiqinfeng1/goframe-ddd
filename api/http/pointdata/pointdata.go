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
}
