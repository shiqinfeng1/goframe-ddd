package pubsub

import (
	"context"
)

type Message struct {
	Topic    string
	Value    []byte
	MetaData any

	Committer
}

func NewMessage(ctx context.Context) *Message {
	return &Message{}
}
