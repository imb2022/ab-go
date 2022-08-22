package ab

import (
	"time"
)

// KafkaConsumer kafka consumer
type KafkaConsumer struct {
	Enable  bool     `json:"enable" mapstructure:"enable"`
	Debug   bool     `json:"debug" mapstructure:"debug"`
	Errors  bool     `json:"errors" mapstructure:"errors"`
	Version string   `json:"version" mapstructure:"version"`
	Brokers []string `json:"brokers" mapstructure:"brokers"`
	// sarama.BalanceStrategyRange, sarama.BalanceStrategyRoundRobin, sarama.BalanceStrategySticky
	// range, roundrobin, sticky
	Topic string `json:"topic" mapstructure:"topic"`
	// Initial Should be OffsetNewest -1 or OffsetOldest -2. Defaults to OffsetNewest.
	Initial int64 `json:"initial" mapstructure:"initial"`
	// ReadUncommitted=0, ReadCommitted=1
	IsolationLevel    int           `json:"isolation_level" mapstructure:"isolation_level"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval" mapstructure:"heartbeat_interval"`
	MaxProcessingTime time.Duration `json:"max_processing_time" mapstructure:"max_processing_time"`
	MaxWaitTime       time.Duration `json:"max_wait_time" mapstructure:"max_wait_time"`
	SessionTimout     time.Duration `json:"session_timeout" mapstructure:"session_timeout"`
}
