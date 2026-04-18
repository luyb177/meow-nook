package config

import (
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Logger logger.Config
}
