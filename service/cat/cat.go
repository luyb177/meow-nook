package main

import (
	"flag"
	"fmt"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/config"
	"github.com/luyb177/meow-nook/service/cat/internal/server"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/cat.yaml", "config file path")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := logger.Init(logger.Config{
		Level:    c.Log.Level,
		Encoding: c.Log.Encoding,
	}); err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}
	defer logger.Sync()

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		server.RegisterCatServer(grpcServer, ctx)
		reflection.Register(grpcServer)
	})
	defer s.Stop()

	logger.Info("cat rpc service starting on " + c.ListenOn)
	s.Start()
}
