package main

import (
	"flag"

	"github.com/luyb177/meow-nook/common/logger"
	grpcmw "github.com/luyb177/meow-nook/common/middleware/grpc"
	"github.com/luyb177/meow-nook/service/cat/internal/config"
	"github.com/luyb177/meow-nook/service/cat/internal/server"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/internal/worker"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"go.uber.org/zap"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "service/cat/etc/cat.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	if ctx.KafkaProducer != nil {
		defer ctx.KafkaProducer.Close()
	}

	err := logger.Init(c.Logger)
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	rpcServer := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		v1.RegisterCatServiceServer(grpcServer, server.NewCatServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer rpcServer.Stop()

	rpcServer.AddUnaryInterceptors(
		grpcmw.ErrorUnaryServer(),
		grpcmw.AccessLogUnary(),
	)

	var sg service.ServiceGroup
	workers := worker.BuildKafkaWorkers(c, ctx)

	sg.Add(rpcServer)
	sg.Add(workers.Pending)
	sg.Add(workers.Scheduler)
	sg.Add(workers.DLQ)

	logger.Info("Starting rpc server...", zap.String("listen_on", c.ListenOn))
	sg.Start()
}
