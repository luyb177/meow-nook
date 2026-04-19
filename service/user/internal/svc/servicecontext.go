package svc

import (
	"github.com/luyb177/meow-nook/common/cache"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/common/mail"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/config"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config        config.Config
	Mailer        *mail.Mailer
	KafkaProducer *kafka.Producer
	RedisClient   *cache.RedisClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	m := mail.NewMailer(mail.EmailConfig{
		From:     c.Email.From,
		Password: c.Email.Password,
		SMTPHost: c.Email.SMTPHost,
		SMTPPort: c.Email.SMTPPort,
	})

	topics := kafka.BuildTopics(c.Kafka.ServiceName)
	producer := kafka.NewProducer(kafka.ProducerConfig{
		Brokers:         c.Kafka.Brokers,
		Topics:          topics,
		DefaultMaxRetry: c.Kafka.DefaultMaxRetry,
		BaseBackoff:     c.Kafka.BaseBackoff,
	})

	var redisClient *cache.RedisClient
	if c.Redis.Addr != "" {
		rc, err := cache.NewRedisClient(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
		if err != nil {
			logger.Warn("failed to connect to Redis for delay queue; retries will fall back to Kafka retry topic",
				zap.String("addr", c.Redis.Addr),
				zap.Error(err),
			)
		} else {
			redisClient = rc
		}
	}

	return &ServiceContext{
		Config:        c,
		Mailer:        m,
		KafkaProducer: producer,
		RedisClient:   redisClient,
	}
}
