package natsclient

import (
	"context"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// topicName    string
// streamName   string
// consumerName string
type SubsKey []string

func (sk SubsKey) TopicName() string {
	return sk[0]
}
func (sk SubsKey) StreamName() string {
	return sk[1]
}
func (sk SubsKey) ConsumerName() string {
	return sk[2]
}
func (sk SubsKey) String() string {
	return strings.Join(sk, "")
}
func NewSubsKey(t, s, c string) SubsKey {
	return []string{t, s, c}
}

type ConsumeFunc func(ctx context.Context, msg *jetstream.Msg) ([]byte, error)
type SubscribeFunc func(ctx context.Context, msg *nats.Msg) ([]byte, error)

type SubType string

const (
	consumeMessageDelay = 100 * time.Millisecond

	JSNEXT   SubType = "js-next"
	JSFETCH  SubType = "js-fetch"
	SUBASYNC SubType = "sub-async"
	SUBSYNC  SubType = "sub-sync"
)

func (s SubType) IsMsg() bool {
	switch s {
	case SUBASYNC, SUBSYNC:
		return true
	}
	return false
}
func (s SubType) IsStream() bool {
	switch s {
	case JSNEXT, JSFETCH:
		return true
	}
	return false
}
