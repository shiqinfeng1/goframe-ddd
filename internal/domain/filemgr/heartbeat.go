package filemgr

import (
	"context"
)

func heartbeat(ctx context.Context, body []byte, _ Repository) []byte {
	return []byte("heartbeat not implemented")
}
