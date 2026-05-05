// internal/logic/cat/list_pending_applies_logic.go
package cat

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ListPendingAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPendingAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPendingAppliesLogic {
	return &ListPendingAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPendingAppliesLogic) ListPendingApplies(req *types.ListPendingAppliesReq) (*types.ListPendingAppliesResp, error) {
	//if !isAdmin(l.ctx) {
	//	return nil, errorx.ErrPermissionDenied
	//}

	// todo casbin 权限检测

	resp, err := l.svcCtx.CatRPC.ListPendingApplies(l.ctx, &catpb.ListPendingAppliesRequest{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.PendingApplyItem, 0, len(resp.List))
	for _, item := range resp.List {
		items = append(items, types.PendingApplyItem{
			ApplyId:          item.ApplyId,
			Name:             item.Name,
			Gender:           item.Gender,
			AgeStage:         item.AgeStage,
			DiscoveryAddress: item.DiscoveryAddress,
			ApplicantUserId:  item.ApplicantUserId,
			CreatedAt:        item.CreatedAt.AsTime().Format(time.RFC3339),
		})
	}

	return &types.ListPendingAppliesResp{
		List:  items,
		Total: resp.Total,
	}, nil
}
