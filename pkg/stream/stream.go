package stream

import (
	"context"
	"errors"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
	"github.com/xtaci/smux"
)

type (
	RecvStreamHandleFunc func(context.Context, *smux.Session, *smux.Stream) error
	SendStreamHandleFunc func(context.Context, *smux.Stream) error
)

type Stream struct {
	clientSess          *smux.Session // 客户端首次连接后，会话信息会被缓存
	clientSessIsRunning atomic.Bool
	tranport            Transport
}

func New(ctx context.Context, tranport Transport) *Stream {
	s := &Stream{
		tranport: tranport,
	}
	return s
}

// 服务端接收一个数据流，首次接收握手消息时，会先启动服务，每个客户端的连接会被缓存
func (s *Stream) StartupServer(ctx context.Context, addr string, recvHandler RecvStreamHandleFunc) error {
	return s.tranport.NewServer(ctx, addr, func(conn net.Conn) {
		session, err := newSessoinByServer(conn)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		go func() {
			defer session.Close()
			g.Log().Info(ctx, "session ready to accept stream ...")
			for {
				// 等待接收一个stream
				stream, err := session.AcceptStream()
				if err != nil {
					if errors.Is(err, io.ErrClosedPipe) {
						g.Log().Info(ctx, "stream pipe is closed")
						return
					}
					g.Log().Errorf(ctx, "session accept stream fail:%v", err)
					return
				}
				// 收到一个stream， 需要在10s内完成收发数据（一个stream接收端数据是16M）
				stream.SetDeadline(
					time.Now().Add(time.Duration(g.Cfg().MustGet(ctx, "filemgr.streamDeadline").Int()) * time.Second),
				)
				g.Log().Infof(ctx, "accept stream ok. remote:%v -> local:%v stream.id=%v", stream.RemoteAddr(), stream.LocalAddr(), stream.ID())
				// 在协程中处理数据
				go func(s *smux.Stream) {
					if err := recvHandler(ctx, session, s); err != nil {
						g.Log().Error(ctx, err)
						return
					}
					// 注意：接收消息的流在处理完数据之后，不主动关闭，由发起方关闭stream
				}(stream)
			}
		}()
	})
}

func (s *Stream) StartupClient(ctx context.Context, addr string) {
	startup := func() error {
		if s.clientSessIsRunning.Load() {
			return nil
		}
		conn, err := s.tranport.NewClient(ctx, addr)
		if err != nil {
			return err
		}
		session, err := newSessionByClient(conn)
		if err != nil {
			conn.Close()
			return err
		}
		s.clientSess = session
		s.clientSessIsRunning.Store(true)
		// 会话建立成功， 立即主动发起握手
		err = s.OpenStreamByClient(ctx, func(ctx context.Context, stm *smux.Stream) error {
			// 构造握手消息
			bytes, err := filemgr.HandshakeMsgToBytes(ctx)
			if err != nil {
				return err
			}
			// 发送握手消息
			if _, err := stm.Write(bytes); err != nil {
				g.Log().Fatalf(ctx, "handshake req to server fail:%v", err)
			}
			// 接收响应数据
			if err := filemgr.CheckoutHandshakeAckFromBytes(ctx, stm); err != nil {
				g.Log().Fatalf(ctx, "handshake ack from server fail:%v", err)
			}
			return nil
		})
		if err != nil {
			session.Close()
			conn.Close()
			s.clientSessIsRunning.Store(false)
		}
		return nil
	}
	// 定时3秒检查连通性
	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for range ticker.C {
			if err := startup(); err != nil {
				g.Log().Errorf(ctx, "filemgr connect to server fail:%v", err)
				continue
			}
			if s.clientSess.IsClosed() {
				s.clientSessIsRunning.Store(false)
			}
		}
	}()
}

func (s *Stream) OpenStreamByServer(ctx context.Context, session *smux.Session, handler SendStreamHandleFunc) error {
	return nil
}

func (s *Stream) OpenStreamByClient(ctx context.Context, handler SendStreamHandleFunc) error {
	if !s.clientSessIsRunning.Load() {
		return gerror.New("session is not exist")
	}
	stream, err := s.clientSess.OpenStream()
	if err != nil {
		if errors.Is(err, io.ErrClosedPipe) {
			g.Log().Info(ctx, "stream pipe is closed")
			return nil
		}
		return err
	}
	g.Log().Infof(ctx, "open stresm ok. local:%v -> remote:%v stream.id=%v", stream.LocalAddr(), stream.RemoteAddr(), stream.ID())
	// 在协程中处理数据
	go func(s *smux.Stream) {
		if err := handler(ctx, s); err != nil {
			g.Log().Error(ctx, err)
			// 即使返回失败，也需要关闭stream
		}
		if err := s.Close(); err != nil {
			g.Log().Errorf(ctx, "close stream.id=%v fail:%v", s.ID(), err)
			return
		}
		g.Log().Infof(ctx, "close stream ok. remote:%v -> local:%v stream.id=%v", stream.RemoteAddr(), stream.LocalAddr(), stream.ID())
	}(stream)
	return nil
}
