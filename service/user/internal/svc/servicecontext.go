package svc

import "github.com/luyb177/meow-nook/service/user/internal/config"

// ServiceContext holds shared dependencies for the user service.
type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{Config: c}
}
