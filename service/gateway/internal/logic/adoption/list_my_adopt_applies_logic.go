package adoption

import (
	"context"

	. "github.com/luyb177/meow-nook/service/gateway/internal/logic"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ListMyAdoptAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyAdoptAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyAdoptAppliesLogic {
	return &ListMyAdoptAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyAdoptAppliesLogic) ListMyAdoptApplies(req *types.ListMyAdoptAppliesReq) (*types.ListMyAdoptAppliesResp, error) {
	logger.Info("ListMyAdoptAppliesLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo get userID from token
	userID := uint64(1)

	resp, err := l.svcCtx.CatRPC.ListMyAdoptApplies(l.ctx, &catpb.ListMyAdoptAppliesRequest{
		Status:      req.Status,
		Page:        req.Page,
		PageSize:    req.PageSize,
		ApplicantId: userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.AdoptApplyItemVO, 0, len(resp.Items))
	for _, v := range resp.Items {
		items = append(items, types.AdoptApplyItemVO{
			ApplyId:             v.ApplyId,
			CatId:               v.CatId,
			CatName:             v.CatName,
			CatAvatar:           v.CatAvatar,
			CatGender:           v.CatGender,
			CreditScoreRequired: v.CreditScoreRequired,
			ApplyReason:         v.ApplyReason,
			Status:              v.Status,
			RejectReason:        v.RejectReason,
			CreatedAt:           PBTimeToString(v.CreatedAt),
			ExpiresAt:           PBTimeToString(v.ExpiresAt),
		})
	}

	return &types.ListMyAdoptAppliesResp{
		Items:    items,
		Total:    int64(resp.Total),
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}, nil
}
