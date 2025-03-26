package nats

import (
	"time"
)

const batchSize = 128

// Config defines the Client configuration.
type Config struct {
	Server       string // 服务地址
	Stream       StreamConfig
	ConsumerName string // 消费者名称
	MaxWait      time.Duration
}

// StreamConfig holds stream settings for NATS jStream.
type StreamConfig struct {
	Name       string   // 流名称
	Subjects   []string // 流下面的主题列表
	MaxDeliver int
	MaxWait    time.Duration
	MaxBytes   int64
}

// validateConfigs validates the configuration for NATS jStream.
func validateConfigs(conf *Config) error {
	if conf.Server == "" {
		return errServerNotProvided
	}

	if len(conf.Stream.Subjects) == 0 {
		return errSubjectsNotProvided
	}
	if conf.Stream.Name == "" {
		return errStreamNotProvided
	}
	if conf.ConsumerName == "" {
		return errConsumerNotProvided
	}

	return nil
}
