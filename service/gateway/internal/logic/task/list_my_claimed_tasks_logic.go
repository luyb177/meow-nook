package task

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	taskpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ListMyClaimedTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyClaimedTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyClaimedTasksLogic {
	return &ListMyClaimedTasksLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListMyClaimedTasksLogic) ListMyClaimedTasks(req *types.ListMyClaimedTasksReq) (*types.ListMyClaimedTasksResp, error) {
	logger.Info("ListMyClaimedTasksLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.ListMyClaimedTasks(l.ctx, &taskpb.ListMyClaimedTasksRequest{
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
		UserId:   userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.ClaimItemVO, 0, len(resp.Items))
	for _, v := range resp.Items {
		items = append(items, types.ClaimItemVO{
			ClaimId:      v.ClaimId,
			TaskId:       v.TaskId,
			TaskTitle:    v.TaskTitle,
			TaskType:     v.TaskType,
			CatName:      v.CatName,
			Status:       v.Status,
			RewardPoints: v.RewardPoints,
			IsOverdue:    v.IsOverdue,
			ClaimedAt:    logic.PBTimeToString(v.ClaimedAt),
			DeadlineAt:   logic.PBTimeToString(v.DeadlineAt),
		})
	}

	return &types.ListMyClaimedTasksResp{
		Items:    items,
		Total:    resp.Total,
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}, nil
}
