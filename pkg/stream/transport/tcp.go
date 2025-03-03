package transport

import (
	"context"
	"net"
	"sync/atomic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtcp"
)

type TcpTransport struct {
	server          *gtcp.Server
	serverIsRunning atomic.Bool
}

func NewTcpTransport() *TcpTransport {
	return &TcpTransport{}
}

func (t *TcpTransport) NewServer(ctx context.Context, addr string, handler func(net.Conn)) error {
	if !t.serverIsRunning.Load() {
		t.server = gtcp.NewServer(addr, func(conn *gtcp.Conn) {
			handler(conn)
		})
		t.serverIsRunning.Store(true)
		go func() {
			if err := t.server.Run(); err != nil {
				g.Log().Errorf(ctx, "tcp server run: %v", err)
				t.server.Close()
				t.serverIsRunning.Store(false)
			}
		}()
		g.Log().Info(ctx, "tcp serevr is running at %v", addr)
	}
	g.Log().Warning(ctx, "start tcp server fail: tcp serevr is already running")
	return nil
}

func (t *TcpTransport) CloseServer(ctx context.Context) {
	if t.server != nil {
		t.server.Close()
		t.server = nil
	}
}

func (t *TcpTransport) NewClient(ctx context.Context, remoteAddr string) (net.Conn, error) {
	conn, err := gtcp.NewNetConn(remoteAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
