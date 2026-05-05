package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type AbandonTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAbandonTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbandonTaskLogic {
	return &AbandonTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AbandonTaskLogic) AbandonTask(in *v1.AbandonTaskRequest) (*v1.AbandonTaskResponse, error) {
	logger.Info("AbandonTask called")

	if err := l.svcCtx.Repo.Task.AbandonTaskClaim(l.ctx, in.ClaimId, in.UserId, in.Reason); err != nil {
		return nil, errorx.WrapDBUpdate("放弃任务失败", err)
	}

	return &v1.AbandonTaskResponse{Message: "已放弃任务"}, nil
}
