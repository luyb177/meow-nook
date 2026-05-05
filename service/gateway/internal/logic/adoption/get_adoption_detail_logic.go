// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	. "github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdoptionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdoptionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdoptionDetailLogic {
	return &GetAdoptionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdoptionDetailLogic) GetAdoptionDetail(req *types.GetAdoptionDetailReq) (*types.GetAdoptionDetailResp, error) {
	logger.Info("GetAdoptionDetailLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo get userID from token
	userID := uint64(1)

	if req.AdoptionId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "adoption_id is required", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.GetAdoptionDetail(l.ctx, &catpb.GetAdoptionDetailRequest{
		AdoptionId:  req.AdoptionId,
		RequesterId: userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	out := &types.GetAdoptionDetailResp{}
	if resp.Adoption != nil {
		a := resp.Adoption
		out.Adoption = types.AdoptionVO{
			Id:                a.Id,
			CatId:             a.CatId,
			CatName:           a.CatName,
			CatAvatar:         a.CatAvatar,
			AdopterId:         a.AdopterId,
			AdopterName:       a.AdopterName,
			Status:            a.Status,
			AgreementNo:       a.AgreementNo,
			AgreedAt:          PBTimeToString(a.AgreedAt),
			AdoptedAt:         PBTimeToString(a.AdoptedAt),
			HomeVisitAt:       PBTimeToString(a.HomeVisitAt),
			HomeVisitUserId:   a.HomeVisitUserId,
			HomeVisitRemark:   a.HomeVisitRemark,
			VisitOneWeekAt:    PBTimeToString(a.VisitOneWeekAt),
			VisitOneMonthAt:   PBTimeToString(a.VisitOneMonthAt),
			VisitThreeMonthAt: PBTimeToString(a.VisitThreeMonthAt),
			VisitSixMonthAt:   PBTimeToString(a.VisitSixMonthAt),
			IsReturned:        a.IsReturned,
			ReturnReason:      a.ReturnReason,
			ReturnedAt:        PBTimeToString(a.ReturnedAt),
			Note:              a.Note,
			CreatedAt:         PBTimeToString(a.CreatedAt),
			UpdatedAt:         PBTimeToString(a.UpdatedAt),
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
	if resp.Adopter != nil {
		out.Adopter = types.UserBriefVO{
			Id:     resp.Adopter.Id,
			Name:   resp.Adopter.Name,
			Avatar: resp.Adopter.Avatar,
		}
	}
	return out, nil
}
