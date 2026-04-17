package main

import (
	"flag"
	"fmt"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/config"
	"github.com/luyb177/meow-nook/service/cat/internal/server"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/cat.yaml", "the config file")

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
		v1.RegisterCatServiceServer(grpcServer, server.NewCatServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
