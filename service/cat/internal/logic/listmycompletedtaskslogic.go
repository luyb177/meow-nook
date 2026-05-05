package logic

import (
	"context"
	"errors"

	"github.com/luyb177/meow-nook/common/errorx"
	taskrepo "github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListMyCompletedTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyCompletedTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyCompletedTasksLogic {
	return &ListMyCompletedTasksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyCompletedTasksLogic) ListMyCompletedTasks(in *v1.ListMyCompletedTasksRequest) (*v1.ListMyCompletedTasksResponse, error) {
	if in.UserId == 0 {
		return nil, errors.New("user_id is required")
	}

	filter := taskrepo.TaskClaimFilter{
		UserID:   in.UserId,
		Status:   taskrepo.ClaimStatusCompleted,
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	claims, total, err := l.svcCtx.Repo.Task.ListTaskClaims(l.ctx, filter)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询已完成任务列表失败", err)
	}

	items := make([]*v1.ClaimItemVO, 0, len(claims))
	var totalPoints int32

	for _, v := range claims {
		// 查询关联任务
		task, err := l.svcCtx.Repo.Task.GetTaskByID(l.ctx, v.TaskID)
		if err != nil {
			continue
		}

		// 查询关联猫咪
		catInfo, _ := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, task.CatID)

		item := &v1.ClaimItemVO{
			ClaimId:      v.ID,
			TaskId:       v.TaskID,
			TaskTitle:    task.Title,
			TaskType:     task.TaskType,
			CatName:      catInfo.Name,
			Status:       v.Status,
			RewardPoints: task.FinalRewardPoints,
			IsOverdue:    v.IsOverdue,
		}

		if !v.CreatedAt.IsZero() {
			item.ClaimedAt = timestamppb.New(v.CreatedAt)
		}
		if task.DeadlineAt != nil {
			item.DeadlineAt = timestamppb.New(*task.DeadlineAt)
		}

		items = append(items, item)
		totalPoints += task.FinalRewardPoints
	}

	return &v1.ListMyCompletedTasksResponse{
		Items:       items,
		Total:       total,
		TotalPoints: totalPoints,
		Page:        in.Page,
		PageSize:    in.PageSize,
	}, nil
}
