package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/luyb177/meow-nook/common/logger"
	"go.uber.org/zap"
)

const (
	defaultSchedulerInterval = time.Second
	defaultClaimLimit        = 100
)

type DelaySchedulerConfig struct {
	Interval   time.Duration
	ClaimLimit int
}

type DelayScheduler struct {
	cfg      DelaySchedulerConfig
	queue    *DelayQueue
	producer *Producer

	ctx    context.Context
	cancel context.CancelFunc
}

func NewDelayScheduler(cfg DelaySchedulerConfig, queue *DelayQueue, producer *Producer) *DelayScheduler {
	if cfg.Interval <= 0 {
		cfg.Interval = defaultSchedulerInterval
	}
	if cfg.ClaimLimit <= 0 {
		cfg.ClaimLimit = defaultClaimLimit
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &DelayScheduler{
		cfg:      cfg,
		queue:    queue,
		producer: producer,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *DelayScheduler) Start() {
	logger.Info("delay scheduler started")
	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.tick()
		}
	}
}

func (s *DelayScheduler) Stop() {
	logger.Info("delay scheduler stopped")
	s.cancel()
}

func (s *DelayScheduler) Close() error {
	s.cancel()
	return nil
}

func (s *DelayScheduler) tick() {
	tickCtx, tickCancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer tickCancel()

	reclaimed, err := s.queue.ReclaimExpired(tickCtx)
	if err != nil {
		logger.Error("delay scheduler: reclaim expired error", zap.Error(err))
	} else if reclaimed > 0 {
		logger.Info("delay scheduler: reclaimed expired tasks", zap.Int64("count", reclaimed))
	}

	taskIDs, err := s.queue.ClaimDue(tickCtx, s.cfg.ClaimLimit)
	if err != nil {
		logger.Error("delay scheduler: claim due error", zap.Error(err))
		return
	}

	if len(taskIDs) > 0 {
		logger.Info("delay scheduler: claimed due tasks", zap.Int("count", len(taskIDs)))
	}
	for _, taskID := range taskIDs {
		s.dispatch(taskID)
	}
}

func (s *DelayScheduler) dispatch(taskID string) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	env, err := s.queue.LoadEnvelope(ctx, taskID)
	if err != nil {
		if errors.Is(err, ErrEnvelopeMissing) || errors.Is(err, ErrEnvelopeCorrupt) {
			logger.Error("delay scheduler: dropping unrecoverable task",
				zap.String("task_id", taskID),
				zap.Error(err),
			)
			return
		}
		logger.Error("delay scheduler: load envelope error",
			zap.String("task_id", taskID),
			zap.Error(err),
		)
		return
	}

	env.NextRunAt = 0
	if err := s.producer.writeEnvelope(ctx, s.producer.cfg.Topics.Pending, env); err != nil {
		logger.Error("delay scheduler: publish to pending error",
			zap.String("task_id", taskID),
			zap.Error(err),
		)
		return
	}

	if err := s.queue.Ack(ctx, taskID); err != nil {
		logger.Error("delay scheduler: ack error",
			zap.String("task_id", taskID),
			zap.Error(err),
		)
	}
}
