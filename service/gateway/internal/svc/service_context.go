// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"github.com/luyb177/meow-nook/service/gateway/internal/config"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user/v1"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRPC userpb.UserServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	uc := zrpc.MustNewClient(c.UserRPC).Conn()

	return &ServiceContext{
		Config:  c,
		UserRPC: userpb.NewUserServiceClient(uc),
	}
}
