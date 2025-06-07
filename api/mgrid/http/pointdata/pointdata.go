// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package pointdata

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/pointdata/v1"
)

type IPointdataV1 interface {
	GetPointData(ctx context.Context, req *v1.GetPointDataReq) (res *v1.GetPointDataRes, err error)
}
