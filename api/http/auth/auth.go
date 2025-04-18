// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package auth

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/http/auth/v1"
)

type IAuthV1 interface {
	DeviceAuth(ctx context.Context, req *v1.DeviceAuthReq) (res *v1.DeviceAuthRes, err error)
}
