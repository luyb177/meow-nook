package config

import "github.com/zeromicro/go-zero/rest"

// Config is the gateway configuration loaded from gateway.yaml.
type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	Casbin struct {
		ModelPath  string
		PolicyPath string
	}

	UserRpc struct {
		Target string
	}

	CatRpc struct {
		Target string
	}

	TaskRpc struct {
		Target string
	}

	AdoptionRpc struct {
		Target string
	}

	PostRpc struct {
		Target string
	}

	Log struct {
		Level    string
		Encoding string
	}
}
