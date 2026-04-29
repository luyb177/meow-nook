package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type EscalateTaskUrgencyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEscalateTaskUrgencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EscalateTaskUrgencyLogic {
	return &EscalateTaskUrgencyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EscalateTaskUrgencyLogic) EscalateTaskUrgency(in *v1.EscalateTaskUrgencyRequest) (*v1.EscalateTaskUrgencyResponse, error) {
	logger.Info("EscalateTaskUrgency called")

	if err := l.svcCtx.Repo.Task.EscalateTaskUrgency(l.ctx, in.TaskId); err != nil {
		return nil, errorx.WrapDBUpdate("升级任务紧急度失败", err)
	}

	task, _ := l.svcCtx.Repo.Task.GetTaskByID(l.ctx, in.TaskId)

	return &v1.EscalateTaskUrgencyResponse{
		NewUrgencyLevel: task.UrgencyLevel,
		EscalationCount: task.EscalationCount,
		Message:         "紧急度已升级",
	}, nil
}
