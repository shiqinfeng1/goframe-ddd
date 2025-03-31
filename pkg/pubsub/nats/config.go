package nats

import (
	"time"
)

// Config defines the Client configuration.
type Config struct {
	Server       string // 服务地址
	ConsumerName string // 消费者名称
	MaxWait      time.Duration
}
