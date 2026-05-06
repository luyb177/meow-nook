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

type ListMyTaskAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyTaskAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyTaskAppliesLogic {
	return &ListMyTaskAppliesLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListMyTaskAppliesLogic) ListMyTaskApplies(req *types.ListMyTaskAppliesReq) (*types.ListMyTaskAppliesResp, error) {
	logger.Info("ListMyTaskAppliesLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.ListMyTaskApplies(l.ctx, &taskpb.ListMyTaskAppliesRequest{
		Status:          req.Status,
		Page:            req.Page,
		PageSize:        req.PageSize,
		ApplicantUserId: uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.TaskApplyItemVO, 0, len(resp.Items))
	for _, v := range resp.Items {
		items = append(items, types.TaskApplyItemVO{
			ApplyId:   v.ApplyId,
			CatId:     v.CatId,
			CatName:   v.CatName,
			CatAvatar: v.CatAvatar,
			Title:     v.Title,
			TaskType:  v.TaskType,
			Status:    v.Status,
			CreatedAt: logic.PBTimeToString(v.CreatedAt),
		})
	}

	return &types.ListMyTaskAppliesResp{
		Items:    items,
		Total:    resp.Total,
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}, nil
}
