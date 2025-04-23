package filemgr

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shiqinfeng1/goframe-ddd/pkg/utils"
)

var MyClientID string

func clientIdFromBytes(_ context.Context, data []byte) string {
	raw := gconv.String(data)
	if !utils.UidIsValid(raw) {
		return ""
	}
	return raw
}

func handshakeBody(_ context.Context) ([]byte, error) {
	uid, err := utils.GenUIDForHost()
	if err != nil {
		return nil, gerror.Wrap(err, "get uid fail")
	}
	MyClientID = uid
	return []byte(uid), nil
}

func handshakeMsgToBytes(ctx context.Context) ([]byte, error) {
	body, err := handshakeBody(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]byte, headerLen+len(body))
	copy(data[0:3], []byte(reqMagic))
	data[3] = msgHandshake.Byte()
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(body)))

	copy(data[8:], body)
	return data, nil
}

func handshakeAckToBytes(_ context.Context, body []byte) ([]byte, error) {
	data := make([]byte, headerLen+len(body))
	copy(data[0:3], []byte(ackMagic))
	data[3] = msgHandshake.Byte()
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(body)))

	copy(data[8:], body)
	return data, nil
}

func recvAck(ctx context.Context, stream io.ReadWriter, mtype msgType) ([]byte, error) {
	header, err := recvHeader(ctx, stream)
	if err != nil {
		return nil, err
	}
	if !header.typ.Is(mtype) {
		return nil, gerror.New("ack fail: msg type not match")
	}
	body, err := recvBody(ctx, stream, header.BodyLen())
	if err != nil {
		return nil, err
	}
	return body, nil
}

func ReqHandshakeWithSync(ctx context.Context, stream io.ReadWriter) error {
	// 构造握手消息
	bytes, err := handshakeMsgToBytes(ctx)
	if err != nil {
		return err
	}
	// 发送握手消息
	if _, err := stream.Write(bytes); err != nil {
		gerror.Wrapf(err, "handshake req to server fail:%v", err)
	}
	// 接收响应数据
	body, err := recvAck(ctx, stream, msgHandshake)
	if err != nil {
		return err
	}
	if clientIdFromBytes(ctx, body) == "" {
		return gerror.New("handshake ack nodeid invalid")
	}
	return nil
}
