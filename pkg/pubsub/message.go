package pubsub

import (
	"encoding/json"

	"github.com/gogf/gf/v2/util/gconv"
)

type Message struct {
	Topic    string
	Value    []byte
	MetaData any
	Subject  string
}

func NewMessage() *Message {
	return &Message{}
}

func (m Message) String() string {
	v := map[string]string{
		"topic":    m.Topic,
		"subject":  m.Subject,
		"value":    gconv.String(m.Value),
		"metaData": gconv.String(m.MetaData),
	}
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
