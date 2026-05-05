package worker

import (
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/cat/internal/config"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/internal/worker/handlers"
)

type KafkaWorkers struct {
	Pending   *kafka.PendingWorker
	Scheduler *kafka.DelayScheduler
	DLQ       *kafka.DLQWatcher
}

func BuildKafkaWorkers(c config.Config, svcCtx *svc.ServiceContext) KafkaWorkers {
	topics := kafka.BuildTopics(c.Kafka.ServiceName)

	pending := kafka.NewPendingWorker(kafka.PendingWorkerConfig{
		Brokers:     c.Kafka.Brokers,
		GroupID:     c.Kafka.GroupID + ".pending",
		Topic:       topics.Pending,
		BaseBackoff: c.Kafka.BaseBackoff,
	}, svcCtx.KafkaProducer)

	dq := kafka.NewDelayQueue(svcCtx.RedisClient.Client, c.Kafka.ServiceName)
	pending.SetDelayQueue(dq)
	scheduler := kafka.NewDelayScheduler(kafka.DelaySchedulerConfig{}, dq, svcCtx.KafkaProducer)

	reg := handlers.BuildRegistry(svcCtx)
	kafka.BindPending(pending, reg)

	dlqNotifier := kafka.NewEmailDLQNotifier(c.Kafka.ServiceName, svcCtx.Mailer, c.DLQEmail.To)

	dlq := kafka.NewDLQWatcher(kafka.DLQWatcherConfig{
		Brokers: c.Kafka.Brokers,
		GroupID: c.Kafka.GroupID + ".dlq",
		Topic:   topics.DLQ,
	}, dlqNotifier)

	return KafkaWorkers{Pending: pending, Scheduler: scheduler, DLQ: dlq}
}
