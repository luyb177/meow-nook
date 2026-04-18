package svc

import (
	"github.com/luyb177/meow-nook/common/cache"
	"github.com/luyb177/meow-nook/common/database"
	"github.com/luyb177/meow-nook/common/mail"
	"github.com/luyb177/meow-nook/service/user/internal/config"
)

// ServiceContext holds shared dependencies for the user service.
type ServiceContext struct {
	Config      config.Config
	MysqlClient *database.MySQLClient
	RedisClient *cache.RedisClient
	Mailer      *mail.Mailer
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

	return &ServiceContext{
		Config:      c,
		MysqlClient: ms,
		RedisClient: r,
		Mailer:      m,
	}
}
