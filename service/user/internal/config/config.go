package config

import (
	"time"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Logger logger.Config

	Kafka KafkaConf

	Email struct {
		From     string
		Password string
		SMTPHost string
		SMTPPort int
	}

	DLQEmail DLQEmailConf
}

// KafkaConf holds Kafka broker and task-queue configuration for the service.
type KafkaConf struct {
	Brokers     []string
	ServiceName string
	GroupID     string

	// DefaultMaxRetry is the default max retry count per task.
	DefaultMaxRetry int
	// BaseBackoff is the base duration for exponential back-off (e.g. "2s").
	BaseBackoff time.Duration
}

// DLQEmailConf holds recipient list for DLQ alert emails.
type DLQEmailConf struct {
	To []string
}
