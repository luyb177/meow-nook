package svc

import (
	"github.com/luyb177/meow-nook/common/cache"
	"github.com/luyb177/meow-nook/common/database"
	"github.com/luyb177/meow-nook/common/mail"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/config"
)

// ServiceContext holds shared dependencies for the user service.
type ServiceContext struct {
	Config        config.Config
	MysqlClient   *database.MySQLClient
	RedisClient   *cache.RedisClient
	Mailer        *mail.Mailer
	KafkaProducer *kafka.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	ms, err := database.NewMySQLClient(c.Mysql.DSN)
	if err != nil {
		panic(err)
	}

	r, err := cache.NewRedisClient(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	if err != nil {
		panic(err)
	}

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

	return &ServiceContext{
		Config:        c,
		MysqlClient:   ms,
		RedisClient:   r,
		Mailer:        m,
		KafkaProducer: producer,
	}
}
