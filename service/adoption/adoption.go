package main

import (
	"flag"
	"fmt"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/adoption/internal/config"
	"github.com/luyb177/meow-nook/service/adoption/internal/server"
	"github.com/luyb177/meow-nook/service/adoption/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/adoption.yaml", "config file path")

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
		server.RegisterAdoptionServer(grpcServer, ctx)
		reflection.Register(grpcServer)
	})
	defer s.Stop()

	logger.Info("adoption rpc service starting on " + c.ListenOn)
	s.Start()
}
