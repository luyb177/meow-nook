package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	DataSource struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	Redis struct {
		Host string
		DB   int
	}

	JWT struct {
		AccessSecret string
		AccessExpire int64
	}

	Log struct {
		Level    string
		Encoding string
	}
}
