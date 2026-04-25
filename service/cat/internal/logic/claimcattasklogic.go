package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClaimCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimCatTaskLogic {
	return &ClaimCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请认领任务
func (l *ClaimCatTaskLogic) ClaimCatTask(in *v1.ClaimCatTaskRequest) (*v1.ClaimCatTaskResponse, error) {
	task, err := l.svcCtx.Repo.Cat.GetCatTaskByID(l.ctx, in.TaskId)
	if err != nil {
		return nil, errorx.ErrTaskNotFound
	}

	if task.Status != "pending" {
		return nil, errorx.ErrTaskAlreadyClaimed
	}

	// Check user doesn't already have an active claim on this task
	_, err = l.svcCtx.Repo.Cat.GetActiveClaimByTaskAndUser(l.ctx, in.TaskId, in.UserId)
	if err == nil {
		return nil, errorx.ErrTaskAlreadyClaimed
	}
	if err != gorm.ErrRecordNotFound {
		logger.Error("ClaimCatTask: check existing claim failed")
		return nil, errorx.WrapInternal("检查认领状态失败", err)
	}

	claim := &cat.CatTaskClaim{
		TaskID:    in.TaskId,
		UserID:    in.UserId,
		Status:    "claimed",
		ClaimTime: time.Now(),
	}
	if err := l.svcCtx.Repo.Cat.CreateCatTaskClaim(l.ctx, claim); err != nil {
		logger.Error("ClaimCatTask: create claim failed")
		return nil, errorx.WrapInternal("认领任务失败", err)
	}

	if err := l.svcCtx.Repo.Cat.UpdateCatTask(l.ctx, in.TaskId, map[string]interface{}{"status": "processing"}); err != nil {
		logger.Error("ClaimCatTask: update task status failed")
		return nil, errorx.WrapInternal("更新任务状态失败", err)
	}

	return &v1.ClaimCatTaskResponse{Success: true}, nil
}
