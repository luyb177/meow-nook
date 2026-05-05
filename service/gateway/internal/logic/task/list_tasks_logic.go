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

type ListTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTasksLogic {
	return &ListTasksLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListTasksLogic) ListTasks(req *types.ListTasksReq) (*types.ListTasksResp, error) {
	logger.Info("ListTasksLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.ListTasks(l.ctx, &taskpb.ListTasksRequest{
		CatId:           req.CatId,
		Status:          req.Status,
		UrgencyLevel:    req.UrgencyLevel,
		DifficultyLevel: req.DifficultyLevel,
		Area:            req.Area,
		Page:            req.Page,
		PageSize:        req.PageSize,
		RequesterId:     userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.TaskItemVO, 0, len(resp.Items))
	for _, v := range resp.Items {
		items = append(items, types.TaskItemVO{
			TaskId:          v.TaskId,
			CatId:           v.CatId,
			CatName:         v.CatName,
			CatAvatar:       v.CatAvatar,
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
			DeadlineAt:      logic.PBTimeToString(v.DeadlineAt),
			CreatedAt:       logic.PBTimeToString(v.CreatedAt),
		})
	}

	return &types.ListTasksResp{
		Items:    items,
		Total:    resp.Total,
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}, nil
}
