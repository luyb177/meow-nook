package svc

import (
	"fmt"

	"github.com/luyb177/meow-nook/common/cache"
	"github.com/luyb177/meow-nook/common/mail"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/config"
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

	// Redis is required for the delay queue (delayed retries).
	if c.Redis.Addr == "" {
		panic("Redis.Addr must be set: Redis is required for delayed retries")
	}
	rc, err := cache.NewRedisClient(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Redis (%s): %v", c.Redis.Addr, err))
	}

	return &ServiceContext{
		Config:        c,
		Mailer:        m,
		KafkaProducer: producer,
		RedisClient:   rc,
	}
}
