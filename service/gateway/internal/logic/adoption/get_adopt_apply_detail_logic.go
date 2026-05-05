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

type GetAdoptApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdoptApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdoptApplyDetailLogic {
	return &GetAdoptApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdoptApplyDetailLogic) GetAdoptApplyDetail(req *types.GetAdoptApplyDetailReq) (*types.GetAdoptApplyDetailResp, error) {
	logger.Info("GetAdoptApplyDetailLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo get userID from token
	userID := uint64(1)

	if req.ApplyId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id is required", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.GetAdoptApplyDetail(l.ctx, &catpb.GetAdoptApplyDetailRequest{
		ApplyId:     req.ApplyId,
		RequesterId: userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	out := &types.GetAdoptApplyDetailResp{}

	if resp.Apply != nil {
		a := resp.Apply
		out.Apply = types.AdoptApplicationVO{
			Id:                   a.Id,
			CatId:                a.CatId,
			CatName:              a.CatName,
			CatAvatar:            a.CatAvatar,
			ApplicantId:          a.ApplicantId,
			ApplicantName:        a.ApplicantName,
			ApplyReason:          a.ApplyReason,
			ContactPhone:         a.ContactPhone,
			ContactWechat:        a.ContactWechat,
			ApplicantCreditScore: a.ApplicantCreditScore,
			Status:               a.Status,
			RejectReason:         a.RejectReason,
			ReviewerId:           a.ReviewerId,
			ReviewerName:         a.ReviewerName,
			ReviewedAt:           PBTimeToString(a.ReviewedAt),
			ApprovedAt:           PBTimeToString(a.ApprovedAt),
			ExpiresAt:            PBTimeToString(a.ExpiresAt),
			CreatedAt:            PBTimeToString(a.CreatedAt),
			UpdatedAt:            PBTimeToString(a.UpdatedAt),
			AdoptionId:           a.AdoptionId,
		}
	}

	if resp.Cat != nil {
		out.Cat = types.CatBriefVO{
			Id:                  resp.Cat.Id,
			Name:                resp.Cat.Name,
			Avatar:              resp.Cat.Avatar,
			Gender:              resp.Cat.Gender,
			CreditScoreRequired: resp.Cat.CreditScoreRequired,
		}
	}

	return out, nil
}
