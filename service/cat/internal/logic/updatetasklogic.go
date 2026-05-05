package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTaskLogic) UpdateTask(in *v1.UpdateTaskRequest) (*v1.UpdateTaskResponse, error) {
	logger.Info("UpdateTask called")

	values := map[string]any{}
	if in.Title != "" {
		values["title"] = in.Title
	}
	if in.Summary != "" {
		values["summary"] = in.Summary
	}
	if in.Detail != "" {
		values["detail"] = in.Detail
	}
	if in.DeadlineAt != nil {
		values["deadline_at"] = in.DeadlineAt.AsTime()
	}
	if in.Remark != "" {
		values["remark"] = in.Remark
	}
	if in.RewardPoints > 0 {
		values["reward_points"] = in.RewardPoints
	}
	if in.DifficultyLevel > 0 {
		values["difficulty_level"] = in.DifficultyLevel
	}

	if err := l.svcCtx.Repo.Task.UpdateTask(l.ctx, in.TaskId, values); err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &v1.UpdateTaskResponse{Message: "更新成功"}, nil
}
