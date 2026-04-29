package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetTaskDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskDetailLogic {
	return &GetTaskDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskDetailLogic) GetTaskDetail(in *v1.GetTaskDetailRequest) (*v1.GetTaskDetailResponse, error) {
	logger.Info("GetTaskDetail called")

	task, err := l.svcCtx.Repo.Task.GetTaskByID(l.ctx, in.TaskId)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询任务失败", err)
	}

	vo := &v1.TaskVO{
		Id:                task.ID,
		CatId:             task.CatID,
		ApplyId:           task.ApplyID,
		CreatorId:         task.CreatorID,
		Title:             task.Title,
		TaskType:          task.TaskType,
		Summary:           task.Summary,
		Detail:            task.Detail,
		Location:          task.Location,
		Longitude:         task.Longitude,
		Latitude:          task.Latitude,
		Area:              task.Area,
		UrgencyLevel:      task.UrgencyLevel,
		DifficultyLevel:   task.DifficultyLevel,
		RewardPoints:      task.RewardPoints,
		FinalRewardPoints: task.FinalRewardPoints,
		MaxClaimers:       task.MaxClaimers,
		CurrentClaimers:   task.CurrentClaimers,
		Status:            task.Status,
		Remark:            task.Remark,
		CreatedAt:         timestamppb.New(task.CreatedAt),
		UpdatedAt:         timestamppb.New(task.UpdatedAt),
	}
	if task.DeadlineAt != nil {
		vo.DeadlineAt = timestamppb.New(*task.DeadlineAt)
	}
	if task.LastEscalatedAt != nil {
		vo.LastEscalatedAt = timestamppb.New(*task.LastEscalatedAt)
	}

	return &v1.GetTaskDetailResponse{Task: vo}, nil
}
