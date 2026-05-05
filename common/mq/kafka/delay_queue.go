package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ErrEnvelopeMissing = errors.New("envelope missing from store")

var ErrEnvelopeCorrupt = errors.New("envelope payload is corrupt")

const (
	visibilityTimeout = 60 * time.Second
)

type DelayQueue struct {
	rdb         *redis.Client
	namespace   string
	retryZ      string
	payloadHash string
	processingZ string
}

func NewDelayQueue(rdb *redis.Client, namespace string) *DelayQueue {
	return &DelayQueue{
		rdb:         rdb,
		namespace:   namespace,
		retryZ:      fmt.Sprintf("%s:retry:z", namespace),
		payloadHash: fmt.Sprintf("%s:retry:payload", namespace),
		processingZ: fmt.Sprintf("%s:retry:processing", namespace),
	}
}

func (q *DelayQueue) Enqueue(ctx context.Context, env *Envelope, delay time.Duration, stage, cause string) error {
	envCopy := *env
	envCopy.Retry++
	envCopy.NextRunAt = time.Now().Add(delay).Unix()
	envCopy.DLQStage = stage
	envCopy.Error = cause

	payload, err := json.Marshal(&envCopy)
	if err != nil {
		return err
	}

	pipe := q.rdb.TxPipeline()
	pipe.HSet(ctx, q.payloadHash, envCopy.TaskID, payload)
	pipe.ZAdd(ctx, q.retryZ, redis.Z{Score: float64(envCopy.NextRunAt), Member: envCopy.TaskID})
	_, err = pipe.Exec(ctx)
	return err
}

var claimDueLua = redis.NewScript(`
local retryZ      = KEYS[1]
local processingZ = KEYS[2]
local now      = tonumber(ARGV[1])
local deadline = tonumber(ARGV[2])
local lim      = tonumber(ARGV[3])

local items = redis.call('ZRANGEBYSCORE', retryZ, '-inf', now, 'LIMIT', 0, lim)
if #items == 0 then
    return {}
end
for _, taskID in ipairs(items) do
    redis.call('ZREM', retryZ, taskID)
    redis.call('ZADD', processingZ, deadline, taskID)
end
return items
`)

func (q *DelayQueue) ClaimDue(ctx context.Context, limit int) ([]string, error) {
	now := time.Now().Unix()
	deadline := now + int64(visibilityTimeout.Seconds())

	result, err := claimDueLua.Run(ctx, q.rdb,
		[]string{q.retryZ, q.processingZ},
		now, deadline, limit,
	).StringSlice()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

var reclaimExpiredLua = redis.NewScript(`
local processingZ = KEYS[1]
local retryZ      = KEYS[2]
local now = tonumber(ARGV[1])

local items = redis.call('ZRANGEBYSCORE', processingZ, '-inf', now)
if #items == 0 then
    return 0
end
for _, taskID in ipairs(items) do
    redis.call('ZREM', processingZ, taskID)
    redis.call('ZADD', retryZ, now, taskID)
end
return #items
`)

func (q *DelayQueue) ReclaimExpired(ctx context.Context) (int64, error) {
	now := time.Now().Unix()

	result, err := reclaimExpiredLua.Run(ctx, q.rdb,
		[]string{q.processingZ, q.retryZ},
		now,
	).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return result, nil
}

func (q *DelayQueue) Ack(ctx context.Context, taskID string) error {
	pipe := q.rdb.TxPipeline()
	pipe.ZRem(ctx, q.processingZ, taskID)
	pipe.HDel(ctx, q.payloadHash, taskID)
	_, err := pipe.Exec(ctx)
	return err
}

func (q *DelayQueue) Drop(ctx context.Context, taskID string) error {
	pipe := q.rdb.TxPipeline()
	pipe.ZRem(ctx, q.processingZ, taskID)
	pipe.HDel(ctx, q.payloadHash, taskID)
	_, err := pipe.Exec(ctx)
	return err
}

func (q *DelayQueue) LoadEnvelope(ctx context.Context, taskID string) (*Envelope, error) {
	data, err := q.rdb.HGet(ctx, q.payloadHash, taskID).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			if dropErr := q.Drop(ctx, taskID); dropErr != nil {
				logger.Error("delay queue: drop after missing envelope failed",
					zap.String("task_id", taskID),
					zap.Error(dropErr),
				)
			}
			return nil, ErrEnvelopeMissing
		}
		return nil, err
	}
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		if dropErr := q.Drop(ctx, taskID); dropErr != nil {
			logger.Error("delay queue: drop after corrupt envelope failed",
				zap.String("task_id", taskID),
				zap.Error(dropErr),
			)
		}
		return nil, ErrEnvelopeCorrupt
	}
	return &env, nil
}
