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

type ListPendingTaskAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPendingTaskAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPendingTaskAppliesLogic {
	return &ListPendingTaskAppliesLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListPendingTaskAppliesLogic) ListPendingTaskApplies(req *types.ListPendingTaskAppliesReq) (*types.ListPendingTaskAppliesResp, error) {
	logger.Info("ListPendingTaskAppliesLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.ListPendingTaskApplies(l.ctx, &taskpb.ListPendingTaskAppliesRequest{
		CatId:    req.CatId,
		Page:     req.Page,
		PageSize: req.PageSize,
		AdminId:  userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.TaskApplyVO, 0, len(resp.Items))
	for _, a := range resp.Items {
		item := types.TaskApplyVO{
			Id:              a.Id,
			CatId:           a.CatId,
			CatName:         a.CatName,
			CatAvatar:       a.CatAvatar,
			ApplicantUserId: a.ApplicantUserId,
			ApplicantName:   a.ApplicantName,
			ApplicantAvatar: a.ApplicantAvatar,
			Title:           a.Title,
			TaskType:        a.TaskType,
			Summary:         a.Summary,
			Detail:          a.Detail,
			Location:        a.Location,
			Longitude:       a.Longitude,
			Latitude:        a.Latitude,
			Status:          a.Status,
			ReviewerId:      a.ReviewerId,
			ReviewerName:    a.ReviewerName,
			ReviewedAt:      logic.PBTimeToString(a.ReviewedAt),
			RejectReason:    a.RejectReason,
			UrgencyLevel:    a.UrgencyLevel,
			DifficultyLevel: a.DifficultyLevel,
			RewardPoints:    a.RewardPoints,
			TaskId:          a.TaskId,
			DeadlineAt:      logic.PBTimeToString(a.DeadlineAt),
			CreatedAt:       logic.PBTimeToString(a.CreatedAt),
			UpdatedAt:       logic.PBTimeToString(a.UpdatedAt),
		}

		tags := make([]types.TagVO, 0, len(a.Tags))
		for _, t := range a.Tags {
			tags = append(tags, types.TagVO{
				Id:    t.Id,
				Name:  t.Name,
				Type:  t.Type,
				Theme: t.Theme,
			})
		}
		item.Tags = tags
		items = append(items, item)
	}

	return &types.ListPendingTaskAppliesResp{
		Items:    items,
		Total:    resp.Total,
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}, nil
}
