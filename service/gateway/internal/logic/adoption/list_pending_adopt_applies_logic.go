package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	. "github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ListPendingAdoptAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPendingAdoptAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPendingAdoptAppliesLogic {
	return &ListPendingAdoptAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPendingAdoptAppliesLogic) ListPendingAdoptApplies(req *types.ListPendingAdoptAppliesReq) (*types.ListPendingAdoptAppliesResp, error) {
	logger.Info("ListPendingAdoptAppliesLogic called")
	//
	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo get userID from token
	userID := uint64(1)

	resp, err := l.svcCtx.CatRPC.ListPendingAdoptApplies(l.ctx, &catpb.ListPendingAdoptAppliesRequest{
		CatId:    req.CatId,
		Page:     req.Page,
		PageSize: req.PageSize,
		AdminId:  userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.PendingApplyVO, 0, len(resp.Items))
	for _, v := range resp.Items {
		items = append(items, types.PendingApplyVO{
			ApplyId:              v.ApplyId,
			CatId:                v.CatId,
			CatName:              v.CatName,
			CatAvatar:            v.CatAvatar,
			CatGender:            v.CatGender,
			CreditScoreRequired:  v.CreditScoreRequired,
			ApplicantId:          v.ApplicantId,
			ApplicantName:        v.ApplicantName,
			ApplicantAvatar:      v.ApplicantAvatar,
			ApplicantCreditScore: v.ApplicantCreditScore,
			ApplyReason:          v.ApplyReason,
			ContactPhone:         v.ContactPhone,
			ContactWechat:        v.ContactWechat,
			CreatedAt:            PBTimeToString(v.CreatedAt),
		})
	}

	return &types.ListPendingAdoptAppliesResp{
		Items:    items,
		Total:    int64(resp.Total),
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}, nil
}
