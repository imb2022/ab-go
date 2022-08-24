package ab

import (
	"time"
)

type Config struct {
	App   string        `json:"app" mapstructure:"app"`
	Kafka KafkaConsumer `json:"kafka" mapstructure:"kafka"`
	MySql MySql         `json:"mysql" mapstructure:"mysql"`
}

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

// MySql mysql config
type MySql struct {
	Enable             bool          `json:"enable" mapstructure:"enable"`
	Addr               string        `json:"addr" mapstructure:"addr"`
	User               string        `json:"user" mapstructure:"user"`
	Password           string        `json:"password" mapstructure:"password"`
	DBName             string        `json:"db_name" mapstructure:"db_name"`
	Charset            string        `json:"charset" mapstructure:"charset"`
	Collation          string        `json:"collation" mapstructure:"collation"`
	MaxOpenConnections int           `json:"max_open_connections" mapstructure:"max_open_connections"`
	MaxIdleConnections int           `json:"max_idle_connections" mapstructure:"max_idle_connections"`
	Debug              bool          `json:"debug" mapstructure:"debug"`
	TLS                bool          `json:"tls" mapstructure:"tls"`
	ParseTime          bool          `json:"parse_time" mapstructure:"parse_time"`
	ConnMaxLifetime    time.Duration `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	Timeout            time.Duration `json:"timeout" mapstructure:"timeout"`
	ReadTimeout        time.Duration `json:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `json:"write_timeout" mapstructure:"write_timeout"`
	LOC                string        `json:"loc" mapstructure:"loc"`
}
