// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package session

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/api/filexfer/http/session/v1"
)

type ISessionV1 interface {
	NodeList(ctx context.Context, req *v1.NodeListReq) (res *v1.NodeListRes, err error)
}
