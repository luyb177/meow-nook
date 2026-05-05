// internal/logic/cat/get_apply_detail_logic.go
package cat

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type GetApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplyDetailLogic {
	return &GetApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApplyDetailLogic) GetApplyDetail(req *types.GetApplyDetailReq) (*types.GetApplyDetailResp, error) {
	resp, err := l.svcCtx.CatRPC.GetApplyDetail(l.ctx, &catpb.GetApplyDetailRequest{
		ApplyId: req.Id,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.GetApplyDetailResp{
		ApplyId:          resp.ApplyId,
		CatId:            resp.CatId,
		Name:             resp.Name,
		Gender:           resp.Gender,
		BodySize:         resp.BodySize,
		AgeStage:         resp.AgeStage,
		Description:      resp.Description,
		DiscoveryAddress: resp.DiscoveryAddress,
		Longitude:        resp.Longitude,
		Latitude:         resp.Latitude,
		ApplicantUserId:  resp.ApplicantUserId,
		Status:           resp.Status,
		RejectReason:     resp.RejectReason,
		CancelReason:     resp.CancelReason,
		ReviewerId:       resp.ReviewerId,
		Images:           convertImagesFromPB(resp.Images),
		CreatedAt:        resp.CreatedAt.AsTime().Format(time.RFC3339),
		UpdatedAt:        resp.UpdatedAt.AsTime().Format(time.RFC3339),
	}, nil
}
