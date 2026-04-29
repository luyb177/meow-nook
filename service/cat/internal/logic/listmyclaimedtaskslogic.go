package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/task"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListMyClaimedTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyClaimedTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyClaimedTasksLogic {
	return &ListMyClaimedTasksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyClaimedTasksLogic) ListMyClaimedTasks(in *v1.ListMyClaimedTasksRequest) (*v1.ListMyClaimedTasksResponse, error) {
	logger.Info("ListMyClaimedTasks called")

	filter := task.TaskClaimFilter{
		UserID:   in.UserId,
		Status:   in.Status,
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	list, total, err := l.svcCtx.Repo.Task.ListTaskClaims(l.ctx, filter)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询已认领任务列表失败", err)
	}

	items := make([]*v1.ClaimItemVO, 0, len(list))
	for _, v := range list {
		item := &v1.ClaimItemVO{
			ClaimId:   v.ID,
			TaskId:    v.TaskID,
			Status:    v.Status,
			IsOverdue: v.IsOverdue,
			ClaimedAt: timestamppb.New(v.CreatedAt),
		}
		if v.CompletedAt != nil {
			item.DeadlineAt = timestamppb.New(*v.CompletedAt)
		}
		items = append(items, item)
	}

	return &v1.ListMyClaimedTasksResponse{
		Items:    items,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
