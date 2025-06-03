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

type ConsumeFunc func(ctx context.Context, msg *jetstream.Msg) error
type SubscribeFunc func(ctx context.Context, msg *nats.Msg) error

type SubType string

const (
	consumeMessageDelay = 100 * time.Millisecond

	SubTypeJSConsumeNext  SubType = "js-next"
	SubTypeJSConsumeFetch SubType = "js-fetch"
	SubTypeSubAsync       SubType = "sub-async"
	SubTypeSubSync        SubType = "sub-sync"
)

func (s SubType) IsMsg() bool {
	switch s {
	case SubTypeSubAsync, SubTypeSubSync:
		return true
	}
	return false
}
func (s SubType) IsStream() bool {
	switch s {
	case SubTypeJSConsumeNext, SubTypeJSConsumeFetch:
		return true
	}
	return false
}
