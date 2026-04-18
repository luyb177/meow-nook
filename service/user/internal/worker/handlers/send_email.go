// Package handlers contains Kafka task handlers for the user service.
package handlers

import (
	"context"
	"encoding/json"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/task"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"go.uber.org/zap"
)

// SendEmailHandler handles user.send_email tasks.
type SendEmailHandler struct {
	svcCtx *svc.ServiceContext
}

func NewSendEmailHandler(svcCtx *svc.ServiceContext) *SendEmailHandler {
	return &SendEmailHandler{svcCtx: svcCtx}
}

func (h *SendEmailHandler) Handle(ctx context.Context, env *kafka.Envelope) error {
	var t task.SendEmailTask
	if err := json.Unmarshal(env.Data, &t); err != nil {
		return err
	}

	logger.Info("handling send_email task",
		zap.String("task_id", env.TaskID),
		zap.String("to", t.To),
		zap.String("subject", t.Subject),
	)

	return h.svcCtx.Mailer.Send(t.Subject, t.Body, []string{t.To})
}
