package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListPendingTaskAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPendingTaskAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPendingTaskAppliesLogic {
	return &ListPendingTaskAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPendingTaskAppliesLogic) ListPendingTaskApplies(in *v1.ListPendingTaskAppliesRequest) (*v1.ListPendingTaskAppliesResponse, error) {
	filter := task.TaskApplyFilter{
		Status:   task.TaskApplyStatusPending,
		CatID:    in.CatId,
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	applies, total, err := l.svcCtx.Repo.Task.ListTaskApplies(l.ctx, filter)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询待审核任务申请列表失败", err)
	}

	// todo it would throw an error if the cat is not found
	items := make([]*v1.TaskApplyVO, 0, len(applies))
	for _, v := range applies {
		// 查询猫咪信息
		catInfo, _ := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, v.CatID)

		item := &v1.TaskApplyVO{
			Id:              v.ID,
			CatId:           v.CatID,
			CatName:         catInfo.Name,
			CatAvatar:       catInfo.Avatar,
			ApplicantUserId: v.ApplicantUserID,
			Title:           v.Title,
			TaskType:        v.TaskType,
			Summary:         v.Summary,
			Detail:          v.Detail,
			Location:        v.Location,
			Longitude:       v.Longitude,
			Latitude:        v.Latitude,
			Status:          v.Status,
			ReviewerId:      v.ReviewerID,
			RejectReason:    v.RejectReason,
			UrgencyLevel:    v.UrgencyLevel,
			DifficultyLevel: v.DifficultyLevel,
			RewardPoints:    v.RewardPoints,
			TaskId:          v.TaskID,
		}

		if v.DeadlineAt != nil {
			item.DeadlineAt = timestamppb.New(*v.DeadlineAt)
		}
		if !v.CreatedAt.IsZero() {
			item.CreatedAt = timestamppb.New(v.CreatedAt)
		}
		if !v.UpdatedAt.IsZero() {
			item.UpdatedAt = timestamppb.New(v.UpdatedAt)
		}

		items = append(items, item)
	}

	return &v1.ListPendingTaskAppliesResponse{
		Items:    items,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
