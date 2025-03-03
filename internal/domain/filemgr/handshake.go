package filemgr

import (
	"context"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/rs/xid"
)

func handshake(ctx context.Context, data []byte) string {
	id, err := xid.FromString(gconv.String(data))
	if err != nil {
		return ""
	}
	return id.String()
}
