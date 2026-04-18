package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/common/mail"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/config"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/task"
	"github.com/luyb177/meow-nook/service/user/internal/server"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/internal/worker/handlers"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

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

	// Build Kafka topics from service name.
	topics := kafka.BuildTopics(c.Kafka.ServiceName)

	// Pending worker – consumes the pending topic and routes by Envelope.Type.
	pendingWorker := kafka.NewPendingWorker(kafka.PendingWorkerConfig{
		Brokers:     c.Kafka.Brokers,
		GroupID:     c.Kafka.GroupID + ".pending",
		Topic:       topics.Pending,
		BaseBackoff: c.Kafka.BaseBackoff,
	}, ctx.KafkaProducer)
	pendingWorker.RegisterHandler(task.TypeSendEmail, handlers.NewSendEmailHandler(ctx))

	// Retry mover – moves expired retry tasks back to pending.
	retryMover := kafka.NewRetryMover(kafka.RetryMoverConfig{
		Brokers:          c.Kafka.Brokers,
		GroupID:          c.Kafka.GroupID + ".retry",
		Topic:            topics.Retry,
		SleepGranularity: time.Second,
	}, ctx.KafkaProducer)

	// DLQ watcher – logs DLQ messages and sends SMTP alerts.
	dlqMailer := mail.NewMailer(mail.EmailConfig{
		From:     c.DLQEmail.From,
		Password: c.DLQEmail.Password,
		SMTPHost: c.DLQEmail.SMTPHost,
		SMTPPort: c.DLQEmail.SMTPPort,
	})
	dlqNotifier := kafka.NewEmailDLQNotifier(c.Kafka.ServiceName, dlqMailer, c.DLQEmail.To)
	dlqWatcher := kafka.NewDLQWatcher(kafka.DLQWatcherConfig{
		Brokers: c.Kafka.Brokers,
		GroupID: c.Kafka.GroupID + ".dlq",
		Topic:   topics.DLQ,
	}, dlqNotifier)

	// gRPC server.
	rpcServer := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		v1.RegisterUserServiceServer(grpcServer, server.NewUserServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	var sg service.ServiceGroup
	sg.Add(rpcServer)
	sg.Add(pendingWorker)
	sg.Add(retryMover)
	sg.Add(dlqWatcher)

	defer sg.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	sg.Start()
}
