package filemgr

import (
	"context"
	"encoding/binary"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shiqinfeng1/goframe-ddd/pkg/utils"
	"github.com/xtaci/smux"
)

var MyClientID string

func clientIdFromBytes(ctx context.Context, data []byte) string {
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

func HandshakeMsgToBytes(ctx context.Context) ([]byte, error) {
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

func HandshakeAckToBytes(ctx context.Context, body []byte) ([]byte, error) {
	data := make([]byte, headerLen+len(body))
	copy(data[0:3], []byte(ackMagic))
	data[3] = msgHandshake.Byte()
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(body)))

	copy(data[8:], body)
	return data, nil
}

func CheckoutHandshakeAckFromBytes(ctx context.Context, stream *smux.Stream) error {
	header, err := recvHeader(ctx, stream)
	if err != nil {
		return err
	}
	if !header.typ.Is(msgHandshake) {
		return gerror.New("not handshake ack")
	}
	body, err := recvBody(ctx, stream, header.BodyLen())
	if err != nil {
		return err
	}
	if clientIdFromBytes(ctx, body) == "" {
		return gerror.New("handshake ack clientid invalid")
	}
	return nil
}
