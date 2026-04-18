package main

import (
	"flag"

	"github.com/luyb177/meow-nook/common/logger"
	grpcmw "github.com/luyb177/meow-nook/common/middleware/grpc"
	"github.com/luyb177/meow-nook/service/user/internal/config"
	"github.com/luyb177/meow-nook/service/user/internal/server"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"
	"go.uber.org/zap"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "service/user/etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	err := logger.Init(c.Logger)
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		v1.RegisterUserServiceServer(grpcServer, server.NewUserServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	s.AddUnaryInterceptors(
		grpcmw.AccessLogUnary(),
		grpcmw.ErrorUnaryServer(),
	)

	logger.Info("Starting rpc server...", zap.String("listen_on", c.ListenOn))
	s.Start()
}
