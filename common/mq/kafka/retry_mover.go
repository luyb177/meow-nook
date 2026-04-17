package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type RetryMoverConfig struct {
	Brokers []string
	GroupID string
	Topic   string

	SleepGranularity time.Duration // 例如 1s
}

type RetryMover struct {
	cfg      RetryMoverConfig
	reader   *kafka.Reader
	producer *Producer

	ctx    context.Context
	cancel context.CancelFunc
}

func NewRetryMover(cfg RetryMoverConfig, producer *Producer) *RetryMover {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MaxWait:  1 * time.Second,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	ctx, cancel := context.WithCancel(context.Background())
	return &RetryMover{
		cfg:      cfg,
		reader:   r,
		producer: producer,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (m *RetryMover) Start() {
	logger.Info("start retry mover")
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			msg, err := m.reader.FetchMessage(m.ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				logger.Error("fetch retry message error", zap.Error(err))
				continue
			}

			var env Envelope
			if err := json.Unmarshal(msg.Value, &env); err != nil {
				// retry 消息坏了：直接 commit 丢弃，并打日志，投 DLQ
				logger.Error("unmarshal retry message error", zap.Error(err))
				_ = m.producer.DispatchDLQ(m.ctx, &env, "unmarshal", "retry", err.Error())
				_ = m.reader.CommitMessages(m.ctx, msg)
				continue
			}

			// 未到期：sleep 一会再继续（避免忙等，且 Stop 可打断）
			now := time.Now().Unix()
			if env.NextRunAt > now {
				delta := time.Duration(env.NextRunAt-now) * time.Second
				sleep := delta
				if m.cfg.SleepGranularity > 0 && sleep > m.cfg.SleepGranularity {
					sleep = m.cfg.SleepGranularity
				}
				timer := time.NewTimer(sleep)
				select {
				case <-m.ctx.Done():
					timer.Stop()
					return
				case <-timer.C:
				}
				// 注意：这里不 commit、不推进 offset，会再次拿到同一条
				continue
			}

			// 到期：重新投递到 pending
			env.NextRunAt = 0
			if err := m.producer.writeEnvelope(m.ctx, m.producer.cfg.Topics.Pending, &env); err != nil {
				logger.Error("dispatch pending from retry error", zap.Error(err))
				// 不 commit，让它下次再试（避免丢）
				continue
			}

			_ = m.reader.CommitMessages(m.ctx, msg)
		}
	}
}

func (m *RetryMover) Stop() {
	logger.Info("stopping retry mover")
	m.cancel()
}

func (m *RetryMover) Close() error {
	m.cancel()
	return m.reader.Close()
}
