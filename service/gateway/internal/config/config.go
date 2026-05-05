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
	CatRPC  zrpc.RpcClientConf

	Logger logger.Config

	JWT JWTConf

	AIService AIServiceConf `json:"AI_Service" yaml:"AI_Service"`
}

// JWTConf holds JWT verification configuration for the gateway.
type JWTConf struct {
	Secret     string        // must match the secret used by user service
	ExpireTime time.Duration // informational only; actual expiry is in the token
}

type AIServiceConf struct {
	ModelAPIKey  string `json:"MODEL_API_KEY" yaml:"MODEL_API_KEY"`
	ModelBaseURL string `json:"MODEL_BASE_URL" yaml:"MODEL_BASE_URL"`
	ModelName    string `json:"MODEL_NAME" yaml:"MODEL_NAME"`
}
