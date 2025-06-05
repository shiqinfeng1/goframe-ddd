package service

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shiqinfeng1/goframe-ddd/pkg/session"
	"github.com/xtaci/smux"
)

func (s *StreamMgr) startupClient(ctx context.Context) error {
	if s.clientSessIsRunning.Load() {
		return nil
	}
	conn, err := s.transport.NewClient(ctx, s.addr)
	if err != nil {
		return err
	}
	sess, err := session.NewSessionByClient(conn)
	if err != nil {
		conn.Close()
		return err
	}
	s.clientConn = conn
	s.clientSess = sess
	s.clientSessIsRunning.Store(true)
	// 会话建立成功， 立即主动发起握手
	err = s.SendByClient(ctx, func(stm *smux.Stream) error {
		c := gctx.New()
		if err := s.ReqHandshake(c, stm); err != nil {
			return gerror.Wrap(err, "handshake fail")
		}
		g.Log().Infof(c, "handshake to server:%v ok", s.addr)
		return nil
	})
	if err != nil {
		g.Log().Warningf(ctx, "handshake to server fail:%v", err)
		sess.Close()
		conn.Close()
		s.clientSessIsRunning.Store(false)
	}
	go s.acceptStream(ctx, sess)
	return err
}

func (s *StreamMgr) StartupClient(ctx context.Context) {
	// 定时3秒检查连通性
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			if err := s.startupClient(ctx); err != nil {
				g.Log().Errorf(ctx, "file mgr connect to server fail:%v", err)
				continue
			}
			if s.clientSess.IsClosed() {
				s.clientSessIsRunning.Store(false)
			}
		}
	}()
}

// 实现文件发送接口
func (s *StreamMgr) SendByClient(ctx context.Context, handler SendStreamHandleFunc) error {
	if s.IsCloud {
		return gerror.New("my is server, cannot send by client")
	}
	if !s.clientSessIsRunning.Load() {
		return gerror.New("session is not exist")
	}
	return s.send(ctx, s.clientSess, handler)
}
