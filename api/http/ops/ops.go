// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package ops

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/http/ops/v1"
)

type IOpsV1 interface {
	Upgrade(ctx context.Context, req *v1.UpgradeReq) (res *v1.UpgradeRes, err error)
}
