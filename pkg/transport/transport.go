package transport

import (
	"context"
	"net"
)

type Transport interface {
	NewClient(ctx context.Context, remoteAddr string) (net.Conn, error)
	NewServer(ctx context.Context, addr string, handler func(net.Conn)) error
	CloseServer(ctx context.Context)
}
