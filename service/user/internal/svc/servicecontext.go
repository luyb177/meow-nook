package svc

import (
	"fmt"

	"github.com/luyb177/meow-nook/common/cache"
	"github.com/luyb177/meow-nook/common/database"
	"github.com/luyb177/meow-nook/common/mail"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/config"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/email"
	"github.com/luyb177/meow-nook/service/user/internal/repo"
)

type ServiceContext struct {
	Config        config.Config
	Mailer        *mail.Mailer // 暂时保留用于构建 kafka 的 dlq
	EmailSender   email.EmailSender
	KafkaProducer *kafka.Producer
	RedisClient   *cache.RedisClient // 同样保留用于构建 kafka 的延时队列
	Repo          *repo.Repositories
}

func NewServiceContext(c config.Config) *ServiceContext {
	m := mail.NewMailer(mail.EmailConfig{
		From:     c.Email.From,
		Password: c.Email.Password,
		SMTPHost: c.Email.SMTPHost,
		SMTPPort: c.Email.SMTPPort,
	})

	emailSender := email.NewEmailSender(m)

	topics := kafka.BuildTopics(c.Kafka.ServiceName)
	producer := kafka.NewProducer(kafka.ProducerConfig{
		Brokers:         c.Kafka.Brokers,
		Topics:          topics,
		DefaultMaxRetry: c.Kafka.DefaultMaxRetry,
		BaseBackoff:     c.Kafka.BaseBackoff,
	})

	if c.RedisConf.Addr == "" {
		panic("RedisConf.Addr must be set: RedisConf is required for delayed retries")
	}
	rc, err := cache.NewRedisClient(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.DB)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Redis (%s): %v", c.RedisConf.Addr, err))
	}

	if c.MySQLConf.DSN == "" {
		panic("MySQL.DSN must be set")
	}
	mysqlClient, err := database.NewMySQLClient(c.MySQLConf.DSN)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to MySQL: %v", err))
	}
	return &ServiceContext{
		Config:        c,
		Mailer:        m,
		EmailSender:   emailSender,
		KafkaProducer: producer,
		RedisClient:   rc,
		Repo:          repo.NewRepositoriesWithDB(rc.Client, mysqlClient.DB),
	}
}
