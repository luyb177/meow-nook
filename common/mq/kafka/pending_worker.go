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

	handlers   map[string]TaskHandler // type string -> handler
	delayQueue *DelayQueue            // optional; if set, retries are scheduled via Redis

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

// SetDelayQueue attaches a Redis-backed delay queue so that retries are
// scheduled there instead of being written back to the Kafka retry topic.
func (w *PendingWorker) SetDelayQueue(q *DelayQueue) {
	w.delayQueue = q
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

			h := w.findHandler(env.Type)
			if h == nil {
				_ = w.producer.DispatchDLQ(w.ctx, &env, "no_handler", "pending", "no handler for task type: "+env.Type)
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

				delay := ExponentialBackoff(w.cfg.BaseBackoff, env.Retry+1)

				if w.delayQueue != nil {
					// Schedule via Redis delay queue.
					if schedErr := w.delayQueue.Enqueue(w.ctx, &env, delay, "handle", err.Error()); schedErr != nil {
						logger.Error("schedule retry via delay queue failed, not committing to preserve message",
							zap.String("task_id", env.TaskID),
							zap.Error(schedErr),
						)
						// Conservative: do not commit so the message is re-processed.
						continue
					}
				} else {
					// Fallback: write to Kafka retry topic.
					_ = w.producer.DispatchRetry(w.ctx, &env, delay, "handle", err.Error())
				}

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
	if err := w.reader.Close(); err != nil {
		logger.Error("close reader close error", zap.Error(err))
	}
}

// findHandler looks up a handler by the exact task type string (e.g. "user.send_email").
func (w *PendingWorker) findHandler(taskType string) TaskHandler {
	return w.handlers[taskType]
}
