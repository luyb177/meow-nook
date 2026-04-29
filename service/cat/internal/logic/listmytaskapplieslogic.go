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

type ListMyTaskAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyTaskAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyTaskAppliesLogic {
	return &ListMyTaskAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyTaskAppliesLogic) ListMyTaskApplies(in *v1.ListMyTaskAppliesRequest) (*v1.ListMyTaskAppliesResponse, error) {
	logger.Info("ListMyTaskApplies called")

	filter := task.TaskApplyFilter{
		ApplicantUserID: in.ApplicantUserId,
		Status:          in.Status,
		Page:            int(in.Page),
		PageSize:        int(in.PageSize),
	}

	list, total, err := l.svcCtx.Repo.Task.ListTaskApplies(l.ctx, filter)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询任务申请列表失败", err)
	}

	items := make([]*v1.TaskApplyItemVO, 0, len(list))
	for _, v := range list {
		items = append(items, &v1.TaskApplyItemVO{
			ApplyId:      v.ID,
			CatId:        v.CatID,
			Title:        v.Title,
			TaskType:     v.TaskType,
			Status:       v.Status,
			RejectReason: v.RejectReason,
			CreatedAt:    timestamppb.New(v.CreatedAt),
		})
	}

	return &v1.ListMyTaskAppliesResponse{
		Items:    items,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
