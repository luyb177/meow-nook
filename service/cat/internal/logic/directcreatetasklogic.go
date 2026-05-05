package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	taskrepo "github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type DirectCreateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDirectCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DirectCreateTaskLogic {
	return &DirectCreateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DirectCreateTaskLogic) DirectCreateTask(in *v1.DirectCreateTaskRequest) (*v1.DirectCreateTaskResponse, error) {
	logger.Info("DirectCreateTask called")

	task := &taskrepo.CatTask{
		CatID:           in.CatId,
		CreatorID:       in.CreatorId,
		Title:           in.Title,
		TaskType:        in.TaskType,
		Summary:         in.Summary,
		Detail:          in.Detail,
		Location:        in.Location,
		Longitude:       in.Longitude,
		Latitude:        in.Latitude,
		Area:            in.Area,
		UrgencyLevel:    in.UrgencyLevel,
		DifficultyLevel: in.DifficultyLevel,
		RewardPoints:    in.RewardPoints,
		MaxClaimers:     in.MaxClaimers,
		Remark:          in.Remark,
	}

	if in.DeadlineAt != nil {
		t := in.DeadlineAt.AsTime()
		task.DeadlineAt = &t
	}

	if err := l.svcCtx.Repo.Task.CreateTask(l.ctx, task); err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &v1.DirectCreateTaskResponse{TaskId: task.ID, Message: "创建成功"}, nil
}
