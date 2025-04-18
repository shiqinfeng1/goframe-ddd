package streammgr

import (
	"context"
	"net"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/pkg/session"
	"github.com/xtaci/smux"
)

// 服务端接收一个数据流，首次接收握手消息时，会先启动服务，每个客户端的连接会被缓存
func (s *StreamMgr) StartupServer(ctx context.Context) error {
	return s.transport.NewServer(ctx, s.addr, func(conn net.Conn) {
		sess, err := session.NewSessoinByServer(conn)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		go s.acceptStream(ctx, sess)
	})
}

// 实现文件发送接口
func (s *StreamMgr) SendByServer(ctx context.Context, session *smux.Session, handler SendStreamHandleFunc) error {
	if !s.IsCloud {
		return gerror.New("my is client, cannot send by server")
	}
	return s.send(ctx, session, handler)
}
