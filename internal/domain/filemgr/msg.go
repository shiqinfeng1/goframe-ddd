package filemgr

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type msgType struct {
	val  int
	desc string
}

func (m msgType) String() string {
	if m.desc != "" {
		return m.desc
	}
	return fmt.Sprintf("undefined msg:%v", m.val)
}

func (m msgType) Byte() byte {
	return gconv.Byte(m.val)
}

func (m msgType) Int() int {
	return m.val
}

func (m msgType) Is(m2 msgType) bool {
	return m == m2
}

type MsgHandleFunc func(context.Context, []byte) error

var (
	headerLen            = 3 + 1 + 4
	maxMsgBodyLen uint32 = 1024 * 1024 * 1024 // 1GB
	reqMagic             = "req"
	ackMagic             = "ack"
	msgMap               = gmap.NewIntStrMap()
	msgHandlerMap        = make(map[msgType]MsgHandleFunc)
	msgHandshake         = newMsgType(1, "文件收发-握手消息")
	msgHeartbeat         = newMsgType(2, "文件收发-心跳消息")
)

func init() {
	msgHandlerMap[msgHeartbeat] = MsgHandleFunc(heartbeat)
}

func newMsgType(val int, desc string) msgType {
	m := msgType{val: val, desc: desc}
	msgMap.Set(val, desc)
	return m
}

type header struct {
	magic  string
	typ    msgType
	length uint32
}

var invalidHeader = func(v int) header {
	return header{magic: "unknow", typ: msgType{val: v, desc: "unknow"}}
}

func (h header) BodyLen() uint32 {
	return h.length
}

func (h header) String() string {
	return fmt.Sprintf("magic:%v type:%v", h.magic, h.typ)
}

func (h header) ErrIfInvalid() error {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{reqMagic, ackMagic})
	if !s.Contains(h.magic) {
		return gerror.Newf("recv header magic invalid(%v)", h.magic)
	}
	if !msgMap.Contains(h.typ.Int()) {
		return gerror.Newf("recv header type invalid(%v)", h.typ)
	}
	if h.length >= maxMsgBodyLen {
		return gerror.Newf("recv header type invalid(%v)", h.typ)
	}
	return nil
}

func newHeaderFromBytes(b []byte) *header {
	if !msgMap.Contains(int(b[3])) {
		iv := invalidHeader(int(b[3]))
		return &iv
	}
	return &header{
		magic:  gconv.String(b[:3]),
		typ:    msgType{val: int(b[3]), desc: msgMap.Get(int(b[3]))},
		length: binary.LittleEndian.Uint32(b[4:8]),
	}
}
