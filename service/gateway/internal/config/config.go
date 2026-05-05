// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"time"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	UserRPC zrpc.RpcClientConf

	Logger logger.Config

	JWT JWTConf
}

// JWTConf holds JWT verification configuration for the gateway.
type JWTConf struct {
	Secret     string        // must match the secret used by user service
	ExpireTime time.Duration // informational only; actual expiry is in the token
}
