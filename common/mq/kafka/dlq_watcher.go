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

type DLQWatcherConfig struct {
	Brokers []string
	GroupID string
	Topic   string
}

type DLQWatcher struct {
	cfg      DLQWatcherConfig
	reader   *kafka.Reader
	notifier Notifier

	ctx    context.Context
	cancel context.CancelFunc
}

func NewDLQWatcher(cfg DLQWatcherConfig, notifier Notifier) *DLQWatcher {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MaxWait:  1 * time.Second,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	ctx, cancel := context.WithCancel(context.Background())
	return &DLQWatcher{
		cfg:      cfg,
		reader:   r,
		notifier: notifier,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (w *DLQWatcher) Start() {
	logger.Info("dlq watcher started")
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			msg, err := w.reader.FetchMessage(w.ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				logger.Error("fetch dlq message error", zap.Error(err))
				continue
			}

			var env Envelope
			if err := json.Unmarshal(msg.Value, &env); err != nil {
				logger.Error("unmarshal dlq envelope error", zap.Error(err))
				_ = w.reader.CommitMessages(w.ctx, msg)
				continue
			}

			// 先日志
			logger.Error(
				"dlq message",
				zap.String("task_id", env.TaskID),
				zap.String("reason", env.DLQReason),
				zap.String("stage", env.DLQStage),
				zap.String("error", env.Error),
				zap.Int("retry", env.Retry),
				zap.Int("max_retry", env.MaxRetry),
			)

			// 再邮件（失败不影响 commit，避免重复轰炸）
			if w.notifier != nil {
				if err := w.notifier.Notify(w.ctx, &env); err != nil {
					logger.Error(
						"dlq notify error",
						zap.Error(err),
					)
				}
			}

			_ = w.reader.CommitMessages(w.ctx, msg)
		}
	}
}

func (w *DLQWatcher) Stop() {
	logger.Info("dlq watcher stopping")
	w.cancel()
	if err := w.reader.Close(); err != nil {
		logger.Error("dlq reader close error", zap.Error(err))
	}
}

func (w *DLQWatcher) Close() error {
	w.cancel()
	return w.reader.Close()
}
