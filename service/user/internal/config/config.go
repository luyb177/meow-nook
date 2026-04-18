package config

import (
	"time"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Log struct {
		Level    string
		Encoding string
	}

	Mysql struct {
		DSN string
	}

	Redis struct {
		Addr     string
		Password string
		DB       int
	}

	JWT struct {
		AccessSecret string
		AccessExpire int64
	}

	Email struct {
		From     string
		Password string
		SMTPHost string
		SMTPPort int
	}

	Kafka KafkaConf

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

// DLQEmailConf holds SMTP configuration and recipient list for DLQ alert emails.
type DLQEmailConf struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort int
	To       []string
}
