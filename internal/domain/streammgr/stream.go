package streammgr

import (
	"context"
	"errors"
	"io"
	"net"
	"sync/atomic"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shiqinfeng1/goframe-ddd/pkg/transport"
	"github.com/xtaci/smux"
)

type (
	RecvStreamHandleFunc func(*smux.Session, io.ReadWriter) error
	SendStreamHandleFunc func(*smux.Stream) error
	ReqHandshakeFunc     func(context.Context, io.ReadWriter) error
)

// 数据流管理
type StreamIntf interface {
	SendByClient(ctx context.Context, handler SendStreamHandleFunc) error
	SendByServer(ctx context.Context, sess *smux.Session, handler SendStreamHandleFunc) error
}

type StreamMgr struct {
	addr                string
	clientConn          net.Conn      // 客户端首次连接后，缓存的会话信息
	clientSess          *smux.Session // 客户端首次连接后，缓存的会话信息
	clientSessIsRunning atomic.Bool
	transport           transport.Transport
	IsCloud             bool
	RecvHandler         RecvStreamHandleFunc
	ReqHandshake        ReqHandshakeFunc
}

// 新建一个流管理
func New() *StreamMgr {
	ctx := gctx.New()
	var stm *StreamMgr

	// 实例化一个流通道管理服务，流通道支持2种传输层：tcp和kcp
	transType := g.Cfg().MustGet(ctx, "sessionmgr.transport").String()
	switch transType {
	case "kcp":
		stm = &StreamMgr{
			transport: transport.NewKcpTransport(),
		}
	case "tcp":
		stm = &StreamMgr{
			transport: transport.NewTcpTransport(),
		}
	default:
		g.Log().Fatalf(ctx, "config of sessionmgr.transport is invalid:%v", transType)
	}
	stm.IsCloud = g.Cfg().MustGet(ctx, "sessionmgr.isCloud").Bool()
	stm.addr = g.Cfg().MustGet(ctx, "sessionmgr.addr").String()

	return stm
}

func (sm *StreamMgr) acceptStream(ctx context.Context, sess *smux.Session) {
	defer sess.Close()
	g.Log().Info(ctx, "session ready to accept stream ...")
	for {
		// 等待接收一个stream
		stm, err := sess.AcceptStream()
		if err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				g.Log().Warning(ctx, "accept stream fail: stream pipe is closed")
				return
			}
			g.Log().Errorf(ctx, "session accept stream fail:%v", err)
			return
		}
		g.Log().Infof(ctx, "accept stream ok. remote:%v -> local:%v stream.id=%v", stm.RemoteAddr(), stm.LocalAddr(), stm.ID())
		// 在协程中处理数据
		go func(s *smux.Stream) {
			for {
				if err := sm.RecvHandler(sess, s); err != nil {
					if gerror.Is(err, io.EOF) {
						g.Log().Infof(ctx, "exit current stream recv handler ok")
						return
					}
					g.Log().Error(ctx, err)
					return
				}
			}
			// 服务端处理完数据后，不需要主动关闭， 等待发起方主动关闭
			// if err := stm.Close(); err != nil {
			// 	g.Log().Errorf(ctx, "close stream.id=%v fail:%v", stm.ID(), err)
			// 	return
			// }
			// g.Log().Infof(ctx, "close stream ok. remote:%v -> local:%v stream.id=%v", stm.RemoteAddr(), stm.LocalAddr(), stream.ID())
		}(stm)
	}
}

func (s *StreamMgr) send(ctx context.Context, session *smux.Session, handler SendStreamHandleFunc) error {
	stream, err := session.OpenStream()
	if err != nil {
		if errors.Is(err, io.ErrClosedPipe) {
			g.Log().Warning(ctx, "open stream fail: stream pipe is closed")
			return nil
		}
		return err
	}
	g.Log().Infof(ctx, "open stream ok. local:%v -> remote:%v stream.id=%v", stream.LocalAddr(), stream.RemoteAddr(), stream.ID())
	// 在协程中处理数据
	go func(s *smux.Stream) {
		if err := handler(s); err != nil {
			g.Log().Warning(ctx, err)
			// 即使返回失败，也需要关闭stream
		}
		if err := s.Close(); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				g.Log().Warningf(ctx, "client-side close stream.id=%v fail:%v", s.ID(), err)
				return
			}
			g.Log().Errorf(ctx, "close stream.id=%v fail:%v", s.ID(), err)
			return
		}
		g.Log().Infof(ctx, "close stream ok. local:%v -> remote:%v stream.id=%v", stream.LocalAddr(), stream.RemoteAddr(), stream.ID())
	}(stream)
	return nil
}
