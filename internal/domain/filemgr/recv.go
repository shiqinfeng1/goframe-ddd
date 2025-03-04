package filemgr

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/xtaci/smux"
)

// 从kcp的stream中接收数据头
func recvHeader(ctx context.Context, stream *smux.Stream) (*header, error) {
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
	return h, nil
}

func recvBody(ctx context.Context, stream *smux.Stream, bodyLen uint32) ([]byte, error) {
	// 读取消息头
	bodyBytes := make([]byte, bodyLen)
	n, err := stream.Read(bodyBytes)
	if err != nil {
		return nil, gerror.Wrap(err, "recv body fail")
	}
	if n != int(bodyLen) {
		return nil, gerror.Newf("recv body length invalid(%v)", n)
	}
	return bodyBytes, nil
}

func StreamRecvHandler(ctx context.Context, sesion *smux.Session, stream *smux.Stream) error {
	header, err := recvHeader(ctx, stream)
	if err != nil {
		return err
	}
	// 首个消息是握手消息，单独处理，缓存session
	if header.typ.Is(msgHandshake) {
		body, err := recvBody(ctx, stream, header.length)
		if err != nil {
			return err
		}
		// 解析clientid
		clientId := clientIdFromBytes(ctx, body)
		if clientId == "" {
			return gerror.Newf("handshake fail: clientId invalid(%v)", gconv.String(body))
		}
		// 回复握手确认消息
		ack, _ := HandshakeAckToBytes(ctx, []byte(clientId))
		if _, err := stream.Write(ack); err != nil {
			return gerror.Wrap(err, "handshake fail")
		}
		// 缓存会话
		if err := saveSession(ctx, clientId, sesion); err != nil {
			return gerror.Wrap(err, "save session fail")
		}
		return nil
	}
	// 其他消息处理
	if handler, ok := msgHandlerMap[header.typ]; ok {
		body, err := recvBody(ctx, stream, header.length)
		if err != nil {
			return err
		}
		if err := handler(ctx, body); err != nil {
			return err
		}
	} else {
		return gerror.Newf("invald msg type:%v", header.typ)
	}

	return nil
}
