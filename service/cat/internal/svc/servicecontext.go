package svc

import "github.com/luyb177/meow-nook/service/cat/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{Config: c}
}
