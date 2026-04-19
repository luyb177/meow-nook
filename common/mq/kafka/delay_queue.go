package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// visibilityTimeout is how long a claimed task stays in processing before
	// being reclaimed as expired (i.e. the scheduler crashed before acking it).
	visibilityTimeout = 60 * time.Second
)

// DelayQueue is a Redis-backed delay queue for scheduling task retries.
//
// Keys (namespaced by service name, e.g. "user.rpc"):
//   - <ns>:retry:z         – ZSET  score=NextRunAt unix-seconds, member=TaskID
//   - <ns>:retry:payload   – HASH  field=TaskID, value=EnvelopeJSON
//   - <ns>:retry:processing – ZSET score=visibilityDeadline unix-seconds, member=TaskID
type DelayQueue struct {
	rdb         *redis.Client
	namespace   string
	retryZ      string // ZSET for scheduled tasks
	payloadHash string // HASH for envelope payloads
	processingZ string // ZSET for in-flight tasks
}

// NewDelayQueue creates a new DelayQueue backed by rdb, keyed under namespace
// (e.g. the Kafka ServiceName such as "user.rpc").
func NewDelayQueue(rdb *redis.Client, namespace string) *DelayQueue {
	return &DelayQueue{
		rdb:         rdb,
		namespace:   namespace,
		retryZ:      fmt.Sprintf("%s:retry:z", namespace),
		payloadHash: fmt.Sprintf("%s:retry:payload", namespace),
		processingZ: fmt.Sprintf("%s:retry:processing", namespace),
	}
}

// Enqueue schedules env for retry after delay.
// It increments Retry, sets NextRunAt, DLQStage and Error on a copy of env,
// then persists the envelope and adds it to the retry ZSET.
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

// claimDueLua atomically moves up to ARGV[3] tasks whose score <= ARGV[1]
// from KEYS[1] (retry:z) into KEYS[2] (processing:z) with score ARGV[2],
// and returns the moved task IDs.
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

// ClaimDue atomically claims up to limit tasks that are due now.
// Returns the list of claimed task IDs.
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

// reclaimExpiredLua moves all tasks in KEYS[1] (processing:z) whose
// visibility deadline has passed back into KEYS[2] (retry:z) as immediately due.
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

// ReclaimExpired moves processing entries whose visibility deadline has elapsed
// back to retry:z so they can be claimed again.
// Returns the number of reclaimed tasks.
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

// Ack removes taskID from the processing ZSET and deletes its payload.
func (q *DelayQueue) Ack(ctx context.Context, taskID string) error {
	pipe := q.rdb.TxPipeline()
	pipe.ZRem(ctx, q.processingZ, taskID)
	pipe.HDel(ctx, q.payloadHash, taskID)
	_, err := pipe.Exec(ctx)
	return err
}

// LoadEnvelope reads the envelope JSON for taskID from the payload hash.
func (q *DelayQueue) LoadEnvelope(ctx context.Context, taskID string) (*Envelope, error) {
	data, err := q.rdb.HGet(ctx, q.payloadHash, taskID).Bytes()
	if err != nil {
		return nil, err
	}
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	return &env, nil
}
