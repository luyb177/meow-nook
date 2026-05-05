package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTasksLogic {
	return &ListTasksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTasksLogic) ListTasks(in *v1.ListTasksRequest) (*v1.ListTasksResponse, error) {
	logger.Info("ListTasks called")

	filter := task.TaskFilter{
		CatID:           in.CatId,
		Status:          in.Status,
		UrgencyLevel:    in.UrgencyLevel,
		DifficultyLevel: in.DifficultyLevel,
		Area:            in.Area,
		Page:            int(in.Page),
		PageSize:        int(in.PageSize),
	}

	list, total, err := l.svcCtx.Repo.Task.ListTasks(l.ctx, filter)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询任务列表失败", err)
	}

	items := make([]*v1.TaskItemVO, 0, len(list))
	for _, v := range list {
		item := &v1.TaskItemVO{
			TaskId:          v.ID,
			CatId:           v.CatID,
			Title:           v.Title,
			TaskType:        v.TaskType,
			Summary:         v.Summary,
			Location:        v.Location,
			UrgencyLevel:    v.UrgencyLevel,
			DifficultyLevel: v.DifficultyLevel,
			RewardPoints:    v.RewardPoints,
			CurrentClaimers: v.CurrentClaimers,
			MaxClaimers:     v.MaxClaimers,
			Status:          v.Status,
			CreatedAt:       timestamppb.New(v.CreatedAt),
		}
		if v.DeadlineAt != nil {
			item.DeadlineAt = timestamppb.New(*v.DeadlineAt)
		}
		items = append(items, item)
	}

	return &v1.ListTasksResponse{
		Items:    items,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
