// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"github.com/casbin/casbin/v2"
	"github.com/luyb177/meow-nook/service/gateway/internal/config"
	"github.com/luyb177/meow-nook/service/gateway/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

// ServiceContext holds all dependencies shared by handler/logic layers.
type ServiceContext struct {
	Config   config.Config
	Enforcer *casbin.Enforcer

	Auth rest.Middleware

	UserRpc      *zrpc.RpcClientConf
	CatRpcConn   *grpc.ClientConn
	TaskRpcConn  *grpc.ClientConn
	AdoptRpcConn *grpc.ClientConn
	PostRpcConn  *grpc.ClientConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	enforcer, err := casbin.NewEnforcer(c.Casbin.ModelPath, c.Casbin.PolicyPath)
	if err != nil {
		panic("failed to load casbin enforcer: " + err.Error())
	}

	catConn, err := grpc.NewClient(c.CatRpc.Target, grpc.WithInsecure()) //nolint:staticcheck
	if err != nil {
		panic("failed to dial cat rpc: " + err.Error())
	}

	taskConn, err := grpc.NewClient(c.TaskRpc.Target, grpc.WithInsecure()) //nolint:staticcheck
	if err != nil {
		panic("failed to dial task rpc: " + err.Error())
	}

	adoptConn, err := grpc.NewClient(c.AdoptionRpc.Target, grpc.WithInsecure()) //nolint:staticcheck
	if err != nil {
		panic("failed to dial adoption rpc: " + err.Error())
	}

	postConn, err := grpc.NewClient(c.PostRpc.Target, grpc.WithInsecure()) //nolint:staticcheck
	if err != nil {
		panic("failed to dial post rpc: " + err.Error())
	}

	return &ServiceContext{
		Config:       c,
		Enforcer:     enforcer,
		Auth:         middleware.NewAuthMiddleware(c.Auth.AccessSecret).Handle,
		CatRpcConn:   catConn,
		TaskRpcConn:  taskConn,
		AdoptRpcConn: adoptConn,
		PostRpcConn:  postConn,
	}
}
