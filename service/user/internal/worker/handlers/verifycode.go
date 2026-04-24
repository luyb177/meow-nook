package handlers

import (
	"context"

	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/task"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"go.uber.org/zap"
)

type VerifyCodeHandler struct {
	svcCtx *svc.ServiceContext
}

func NewVerifyCodeHandler(svcCtx *svc.ServiceContext) *VerifyCodeHandler {
	return &VerifyCodeHandler{
		svcCtx: svcCtx,
	}
}

func (h *VerifyCodeHandler) SendVerifyCode(ctx context.Context, env *kafka.Envelope) error {
	var p task.VerifyCode
	if err := kafka.Decode(env, &p); err != nil {
		return err
	}

	logger.Info("Sending verify_code task", zap.String("task_id", env.TaskID))

	return h.svcCtx.EmailSender.SendVerifyCode(ctx, p.To, p.Code, 5)
}
