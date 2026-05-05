// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/config"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user/v1"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	UserRPC   userpb.UserServiceClient
	JWTSecret string // convenience shortcut used by middleware builders
	Config  config.Config
	UserRPC userpb.UserServiceClient
	CatRPC  catpb.CatServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	uc := zrpc.MustNewClient(c.UserRPC).Conn()
	cc := zrpc.MustNewClient(c.CatRPC).Conn()
	return &ServiceContext{
		Config:    c,
		UserRPC:   userpb.NewUserServiceClient(uc),
		JWTSecret: c.JWT.Secret,
		Config:  c,
		UserRPC: userpb.NewUserServiceClient(uc),
		CatRPC:  catpb.NewCatServiceClient(cc),
	}
}
