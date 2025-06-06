package filemgr

import (
	"context"
	"io"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shiqinfeng1/goframe-ddd/pkg/session"
	"github.com/xtaci/smux"
)

// 从kcp的stream中接收数据头
func recvHeader(_ context.Context, stream io.Reader) (*header, error) {
	// 读取消息头
	headerBytes := make([]byte, headerLen)
	n, err := stream.Read(headerBytes)
	if err != nil {
		return nil, gerror.Wrap(err, "recv header fail")
	}
	if n != headerLen {
		return nil, gerror.Wrapf(err, "recv header length invalid(%v)", n)
	}
	// 解析消息头
	h := newHeaderFromBytes(headerBytes)
	if err := h.ErrIfInvalid(); err != nil {
		return nil, err
	}
	// g.Log().Debugf(ctx, "recv msg header: %v", h)
	return h, nil
}

func recvBody(_ context.Context, stream io.Reader, bodyLen uint32) ([]byte, error) {
	bodyBytes := make([]byte, bodyLen)
	var m int = 0
	for {
		n, err := stream.Read(bodyBytes[m:])
		if err != nil {
			return nil, gerror.Wrap(err, "recv body fail")
		}
		m += n
		// g.Log().Debugf(ctx, "recv msg body %v bytes", n)
		if m == int(bodyLen) {
			break
		}
	}
	return bodyBytes, nil
}

func ackHandshake(ctx context.Context, sesion *smux.Session, stream io.Writer, body []byte) error {
	// 解析clientid
	nodeId := clientIdFromBytes(ctx, body)
	if nodeId == "" {
		return gerror.Newf("handshake fail: nodeId invalid(%v)", gconv.String(body))
	}
	// 回复握手确认消息
	ack, _ := handshakeAckToBytes(ctx, []byte(nodeId))
	if _, err := stream.Write(ack); err != nil {
		return gerror.Wrap(err, "handshake fail")
	}
	// 缓存会话
	if err := session.SaveSession(ctx, nodeId, sesion); err != nil {
		return gerror.Wrap(err, "save session fail")
	}
	g.Log().Infof(ctx, "handshake from:%v ok", nodeId)
	return nil
}

func (f *FileTransferMgr) StreamRecvHandler(sesion *smux.Session, stm io.ReadWriter) error {
	ctx := gctx.New()
	header, err := recvHeader(ctx, stm)
	if err != nil {
		return err
	}
	body, err := recvBody(ctx, stm, header.length)
	if err != nil {
		return err
	}
	// 首个消息是握手消息，单独处理，缓存session
	if header.typ.Is(msgHandshake) {
		if err := ackHandshake(ctx, sesion, stm, body); err != nil {
			return err
		}
		return nil
	}
	// 其他消息处理
	if handler, ok := msgHandlerMap[header.typ]; ok {
		ack := handler(ctx, body, f.repo)
		// 回复握手确认消息
		if _, err := stm.Write(ack); err != nil {
			return gerror.Wrap(err, "write ack fail")
		}
	} else {
		return gerror.Newf("not register handler for msg type:%v", header.typ)
	}

	return nil
}
