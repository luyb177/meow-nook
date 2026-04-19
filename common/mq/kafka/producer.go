package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Brokers []string
	Topics  Topics

	DefaultMaxRetry int
	BaseBackoff     time.Duration
}

type Producer struct {
	cfg    ProducerConfig
	writer *kafka.Writer
}

func NewProducer(cfg ProducerConfig) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Balancer:     &kafka.Hash{}, // key hash，保证同 taskID 同分区
		RequiredAcks: kafka.RequireOne,
		Async:        true,
	}
	return &Producer{cfg: cfg, writer: w}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) Dispatch(ctx context.Context, task Task, opts ...DispatchOption) error {
	payload, err := task.Payload()
	if err != nil {
		return err
	}

	o := &DispatchOptions{}
	for _, f := range opts {
		f(o)
	}
	maxRetry := p.cfg.DefaultMaxRetry
	if o.MaxRetry != nil {
		maxRetry = *o.MaxRetry
	}

	env := &Envelope{
		TaskID:    BuildTaskID(task.Type(), task.Biz()),
		Type:      task.Type(),
		Retry:     0,
		MaxRetry:  maxRetry,
		CreatedAt: time.Now().Unix(),
		NextRunAt: 0,
		Data:      payload,
	}
	return p.writeEnvelope(ctx, p.cfg.Topics.Pending, env)
}

func (p *Producer) DispatchRetry(ctx context.Context, env *Envelope, delay time.Duration, stage, cause string) error {
	envCopy := *env
	envCopy.Retry++
	envCopy.NextRunAt = time.Now().Add(delay).Unix()
	envCopy.DLQStage = stage
	envCopy.Error = cause
	return p.writeEnvelope(ctx, p.cfg.Topics.Retry, &envCopy)
}

func (p *Producer) DispatchDLQ(ctx context.Context, env *Envelope, reason, stage, cause string) error {
	envCopy := *env
	envCopy.DLQReason = reason
	envCopy.DLQStage = stage
	envCopy.Error = cause
	return p.writeEnvelope(ctx, p.cfg.Topics.DLQ, &envCopy)
}

func (p *Producer) writeEnvelope(ctx context.Context, topic string, env *Envelope) error {
	b, err := json.Marshal(env)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(env.TaskID),
		Value: b,
		Time:  time.Now(),
	}
	return p.writer.WriteMessages(ctx, msg)
}
