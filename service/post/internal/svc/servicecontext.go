package svc

import "github.com/luyb177/meow-nook/service/post/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{Config: c}
}
