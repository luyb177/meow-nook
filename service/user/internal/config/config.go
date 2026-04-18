package config

import "github.com/zeromicro/go-zero/zrpc"

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
}
