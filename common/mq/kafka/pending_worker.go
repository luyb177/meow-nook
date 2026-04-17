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

type TaskHandler interface {
	Handle(ctx context.Context, env *Envelope) error
}

type PendingWorkerConfig struct {
	Brokers []string
	GroupID string
	Topic   string

	BaseBackoff time.Duration
}

type PendingWorker struct {
	cfg      PendingWorkerConfig
	reader   *kafka.Reader
	producer *Producer

	handlers map[string]TaskHandler // prefix -> handler

	ctx    context.Context
	cancel context.CancelFunc
}

func NewPendingWorker(cfg PendingWorkerConfig, producer *Producer) *PendingWorker {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MaxWait:  1 * time.Second,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	ctx, cancel := context.WithCancel(context.Background())
	return &PendingWorker{
		cfg:      cfg,
		reader:   r,
		producer: producer,
		handlers: map[string]TaskHandler{},
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (w *PendingWorker) Close() error {
	w.cancel()
	return w.reader.Close()
}

func (w *PendingWorker) RegisterHandler(prefix string, h TaskHandler) {
	w.handlers[prefix] = h
}

func (w *PendingWorker) Start() {
	logger.Info("pending worker started")
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
				logger.Error("fetch message error", zap.Error(err))
				continue
			}

			var env Envelope
			if err := json.Unmarshal(msg.Value, &env); err != nil {
				// 解析失败：直接进 DLQ，然后 commit，避免卡住
				_ = w.producer.DispatchDLQ(w.ctx, &Envelope{
					TaskID:    string(msg.Key),
					CreatedAt: msg.Time.Unix(),
					Data:      msg.Value,
					MaxRetry:  0,
				}, "unmarshal", "pending", err.Error())

				_ = w.reader.CommitMessages(w.ctx, msg)
				continue
			}

			h := w.findHandler(env.TaskID)
			if h == nil {
				_ = w.producer.DispatchDLQ(w.ctx, &env, "no_handler", "pending", "no handler for task id prefix")
				_ = w.reader.CommitMessages(w.ctx, msg)
				continue
			}

			if err := h.Handle(w.ctx, &env); err != nil {
				// 超过 maxRetry -> DLQ
				if env.Retry >= env.MaxRetry {
					_ = w.producer.DispatchDLQ(w.ctx, &env, "max_retry", "handle", err.Error())
					_ = w.reader.CommitMessages(w.ctx, msg)
					continue
				}

				// 写入 retry topic（指数退避）
				delay := ExponentialBackoff(w.cfg.BaseBackoff, env.Retry+1)
				_ = w.producer.DispatchRetry(w.ctx, &env, delay, "handle", err.Error())

				// 关键：commit 当前 pending offset，避免卡 partition
				_ = w.reader.CommitMessages(w.ctx, msg)
				continue
			}

			// 成功：commit
			if err := w.reader.CommitMessages(w.ctx, msg); err != nil {
				logger.Error("commit messages error", zap.Error(err))
			}
		}
	}
}

func (w *PendingWorker) Stop() {
	logger.Info("pending worker stopped")
	w.cancel()
}

func (w *PendingWorker) findHandler(taskID string) TaskHandler {
	for prefix, handler := range w.handlers {
		if len(taskID) >= len(prefix) && taskID[:len(prefix)] == prefix {
			return handler
		}
	}
	return nil
}
