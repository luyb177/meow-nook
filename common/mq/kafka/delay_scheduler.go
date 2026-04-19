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

// DelaySchedulerConfig configures the DelayScheduler.
type DelaySchedulerConfig struct {
	// Interval is how often the scheduler polls for due tasks (default 1s).
	Interval time.Duration
	// ClaimLimit is the maximum number of tasks claimed per tick (default 100).
	ClaimLimit int
}

// DelayScheduler replaces RetryMover.  It periodically:
//  1. Calls ReclaimExpired to recover stale processing entries.
//  2. Calls ClaimDue to atomically claim tasks that are ready to run.
//  3. Publishes each claimed envelope back to the Kafka pending topic.
//  4. Acks on publish success; leaves in processing on failure (reclaim will retry).
type DelayScheduler struct {
	cfg      DelaySchedulerConfig
	queue    *DelayQueue
	producer *Producer

	ctx    context.Context
	cancel context.CancelFunc
}

// NewDelayScheduler creates a DelayScheduler.
// queue must not be nil; producer is used to write to the pending topic.
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

// Start begins the scheduling loop.  It blocks until Stop/Close is called.
// Implements service.Service.
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

// Stop cancels the scheduling loop. Implements service.Service.
func (s *DelayScheduler) Stop() {
	logger.Info("delay scheduler stopped")
	s.cancel()
}

// Close cancels the scheduling loop. Implements io.Closer.
func (s *DelayScheduler) Close() error {
	s.cancel()
	return nil
}

// tick is called once per interval.
func (s *DelayScheduler) tick() {
	// Use a short context for reclaim and claim operations.
	tickCtx, tickCancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer tickCancel()

	// Step 1: reclaim tasks whose visibility timeout has elapsed.
	reclaimed, err := s.queue.ReclaimExpired(tickCtx)
	if err != nil {
		logger.Error("delay scheduler: reclaim expired error", zap.Error(err))
	} else if reclaimed > 0 {
		logger.Info("delay scheduler: reclaimed expired tasks", zap.Int64("count", reclaimed))
	}

	// Step 2: claim tasks that are due now.
	taskIDs, err := s.queue.ClaimDue(tickCtx, s.cfg.ClaimLimit)
	if err != nil {
		logger.Error("delay scheduler: claim due error", zap.Error(err))
		return
	}

	if len(taskIDs) > 0 {
		logger.Info("delay scheduler: claimed due tasks", zap.Int("count", len(taskIDs)))
	}

	// Step 3: publish each claimed task to the pending topic.
	// Each dispatch uses its own bounded context so a slow Kafka write
	// cannot stall the entire tick.
	for _, taskID := range taskIDs {
		s.dispatch(taskID)
	}
}

// dispatch publishes a single task back to the Kafka pending topic.
func (s *DelayScheduler) dispatch(taskID string) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	env, err := s.queue.LoadEnvelope(ctx, taskID)
	if err != nil {
		// Terminal errors: the payload is gone or corrupt.  Drop already cleaned up
		// the processing entry, so just log and move on.
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
		// Leave in processing; ReclaimExpired will pick it up after the visibility timeout.
		return
	}

	// Reset NextRunAt so the pending worker does not treat this as a scheduled task.
	env.NextRunAt = 0
	if err := s.producer.writeEnvelope(ctx, s.producer.cfg.Topics.Pending, env); err != nil {
		logger.Error("delay scheduler: publish to pending error",
			zap.String("task_id", taskID),
			zap.Error(err),
		)
		// Leave in processing; ReclaimExpired will pick it up.
		return
	}

	if err := s.queue.Ack(ctx, taskID); err != nil {
		logger.Error("delay scheduler: ack error",
			zap.String("task_id", taskID),
			zap.Error(err),
		)
	}
}
