package handlers

import (
	"context"
	"errors"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/task"
	"go.uber.org/zap"
)

// 注册不同 type 的 handler

type TestHandlers struct{}

func (h TestHandlers) OnSuccess(ctx context.Context, env *kafka.Envelope) error {
	var p task.TestPayload
	if err := kafka.Decode(env, &p); err != nil {
		return err
	}

	logger.Info("handling success_test",
		zap.String("task_id", env.TaskID),
		zap.String("name", p.Name),
	)
	return nil
}

func (h TestHandlers) OnFail(ctx context.Context, env *kafka.Envelope) error {
	var p task.TestPayload
	if err := kafka.Decode(env, &p); err != nil {
		return err
	}

	logger.Info("handling fail_test",
		zap.String("task_id", env.TaskID),
		zap.String("name", p.Name),
	)
	return errors.New("fail_test")
}
