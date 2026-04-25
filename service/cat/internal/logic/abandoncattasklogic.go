package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type AbandonCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAbandonCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbandonCatTaskLogic {
	return &AbandonCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请放弃任务
func (l *AbandonCatTaskLogic) AbandonCatTask(in *v1.AbandonCatTaskRequest) (*v1.AbandonCatTaskResponse, error) {
	if _, err := l.svcCtx.Repo.Cat.GetCatTaskByID(l.ctx, in.TaskId); err != nil {
		return nil, errorx.ErrTaskNotFound
	}

	claim, err := l.svcCtx.Repo.Cat.GetActiveClaimByTaskAndUser(l.ctx, in.TaskId, in.UserId)
	if err != nil {
		return nil, errorx.ErrTaskNotOwned
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatTaskClaim(l.ctx, claim.ID, map[string]interface{}{
		"status":         "abandoned",
		"abandon_reason": in.Reason,
	}); err != nil {
		logger.Error("AbandonCatTask: update claim failed")
		return nil, errorx.WrapInternal("更新认领状态失败", err)
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatTask(l.ctx, in.TaskId, map[string]interface{}{"status": "pending"}); err != nil {
		logger.Error("AbandonCatTask: update task status failed")
		return nil, errorx.WrapInternal("更新任务状态失败", err)
	}

	return &v1.AbandonCatTaskResponse{Success: true}, nil
}
