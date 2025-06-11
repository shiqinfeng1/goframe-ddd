package natsclient

import (
	"strings"
	"time"
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

const (
	consumeMessageDelay = 100 * time.Millisecond
)
