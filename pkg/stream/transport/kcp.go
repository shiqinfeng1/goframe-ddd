package transport

import (
	"context"
	"net"
	"sync/atomic"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/xtaci/kcp-go"
)

type KcpTransport struct {
	server          *kcp.Listener
	serverIsRunning atomic.Bool
}

func NewKcpTransport() *KcpTransport {
	return &KcpTransport{}
}

func (k *KcpTransport) NewServer(ctx context.Context, addr string, handler func(net.Conn)) error {
	if k.serverIsRunning.Load() {
		g.Log().Warning(ctx, "start kcp server fail: kcp serevr is already running")
		return nil
	}
	listener, err := kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		return gerror.Wrapf(err, "kcp listen fail")
	}
	listener.SetReadBuffer(4 * 1024 * 1024)
	listener.SetWriteBuffer(4 * 1024 * 1024)
	listener.SetDSCP(46)
	k.server = listener
	k.serverIsRunning.Store(true)

	go func() {
		for {
			conn, err := k.server.AcceptKCP()
			if err != nil {
				g.Log().Errorf(ctx, "kcp accept fail:%v", err)
				k.server.Close()
				k.serverIsRunning.Store(false)
				return
			}
			conn.SetStreamMode(true)
			conn.SetWindowSize(4096, 4096)
			conn.SetNoDelay(1, 10, 2, 1)
			conn.SetWriteDelay(false)
			conn.SetDSCP(46)
			conn.SetMtu(1400)
			conn.SetACKNoDelay(false)
			handler(conn)
		}
	}()
	g.Log().Infof(ctx, "kcp serevr is running at %v", addr)

	return nil
}

func (k *KcpTransport) CloseServer(ctx context.Context) {
	if k.server != nil {
		k.server.Close()
		k.server = nil
	}
}

func (k *KcpTransport) NewClient(ctx context.Context, remoteAddr string) (net.Conn, error) {
	sess, err := kcp.DialWithOptions(remoteAddr, nil, 10, 3)
	if err != nil {
		return nil, gerror.Newf("kcp DialWithOptions:%v", err)
	}
	sess.SetStreamMode(true)
	sess.SetWriteDelay(false)
	sess.SetWindowSize(1024, 1024)
	sess.SetReadBuffer(16 * 1024 * 1024)
	sess.SetWriteBuffer(16 * 1024 * 1024)
	sess.SetNoDelay(1, 10, 2, 1)
	sess.SetMtu(1400)
	sess.SetDSCP(46)
	sess.SetACKNoDelay(false)
	return sess, nil
}
