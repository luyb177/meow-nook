package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCatTaskLogic {
	return &CreateCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员直接创建正式任务（无需审核）
func (l *CreateCatTaskLogic) CreateCatTask(in *v1.CreateCatTaskRequest) (*v1.CreateCatTaskResponse, error) {
	// TODO(permission): check that in.OperatorId has admin role

	if _, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, in.CatId); err != nil {
		return nil, errorx.ErrCatNotFound
	}

	task := &cat.CatTask{
		CatID:           in.CatId,
		Title:           in.Title,
		TaskType:        in.TaskType,
		UrgencyLevel:    in.UrgencyLevel,
		DifficultyLevel: in.DifficultyLevel,
		RewardPoints:    in.RewardPoints,
		Status:          "pending",
		Summary:         in.Summary,
		Detail:          in.Detail,
		Deadline:        in.Deadline.AsTime(),
	}

	if err := l.svcCtx.Repo.Cat.CreateCatTask(l.ctx, task); err != nil {
		logger.Error("CreateCatTask: create task failed")
		return nil, errorx.WrapInternal("创建任务失败", err)
	}

	return &v1.CreateCatTaskResponse{TaskId: task.ID}, nil
}
