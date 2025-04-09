// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package ops

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/http/ops/v1"
)

type IOpsV1 interface {
	UpgradeApp(ctx context.Context, req *v1.UpgradeAppReq) (res *v1.UpgradeAppRes, err error)
	UpgradeImage(ctx context.Context, req *v1.UpgradeImageReq) (res *v1.UpgradeImageRes, err error)
}
