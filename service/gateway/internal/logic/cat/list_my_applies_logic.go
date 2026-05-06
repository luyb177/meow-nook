// internal/logic/cat/list_my_applies_logic.go
package cat

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ListMyAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyAppliesLogic {
	return &ListMyAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyAppliesLogic) ListMyApplies(req *types.ListMyAppliesReq) (*types.ListMyAppliesResp, error) {
	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.ListMyApplies(l.ctx, &catpb.ListMyAppliesRequest{
		ApplicantUserId: uint64(userID),
		Status:          req.Status,
		Page:            req.Page,
		PageSize:        req.PageSize,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.ApplyItem, 0, len(resp.List))
	for _, item := range resp.List {
		items = append(items, types.ApplyItem{
			ApplyId:   item.ApplyId,
			Name:      item.Name,
			Status:    item.Status,
			CatId:     item.CatId,
			CreatedAt: item.CreatedAt.AsTime().Format(time.RFC3339),
		})
	}

	return &types.ListMyAppliesResp{
		List:  items,
		Total: resp.Total,
	}, nil
}
