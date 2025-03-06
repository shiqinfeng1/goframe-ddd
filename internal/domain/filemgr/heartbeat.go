package filemgr

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

func heartbeat(ctx context.Context, data []byte) error {
	return gerror.New("heartbeat not implemented")
}
