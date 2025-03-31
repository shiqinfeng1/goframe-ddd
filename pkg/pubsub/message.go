package pubsub

import (
	"fmt"

	"github.com/gogf/gf/v2/util/gconv"
)

type Message struct {
	Topic    string
	Value    []byte
	MetaData any

	Committer
}

func NewMessage() *Message {
	return &Message{}
}

func (m Message) String() string {
	return fmt.Sprintf("topic:%v value:%v metadata:%v", m.Topic, gconv.String(m.Value), gconv.String(m.MetaData))
}
