package cat

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveCreateCatLogic {
	return &ApproveCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApproveCreateCatLogic) ApproveCreateCat(req *types.ApproveCreateCatReq) (*types.ApproveCreateCatResp, error) {
	//if !isAdmin(l.ctx) {
	//	return nil, errorx.ErrPermissionDenied
	//}
	//
	//uid := getUserID(l.ctx)
	//if uid == 0 {
	//	return nil, errorx.ErrUnauthorized
	//}

	// todo 从 JWT 中获取用户 ID
	uid := uint64(1)
	// todo casbin 权限检测

	resp, err := l.svcCtx.CatRPC.ApproveCreateCat(l.ctx, &catpb.ApproveCreateCatRequest{
		ApplyId:                 req.Id,
		OperatorId:              uid,
		CatCode:                 req.CatCode,
		Breed:                   req.Breed,
		Color:                   req.Color,
		IsVaccinated:            req.IsVaccinated,
		IsHealthy:               req.IsHealthy,
		NeedMedicalIntervention: req.NeedMedicalIntervention,
		SterilizationStatus:     req.SterilizationStatus,
		ExtraTagIds:             req.ExtraTagIds,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.ApproveCreateCatResp{
		CatId:   resp.CatId,
		CatCode: resp.CatCode,
		Status:  resp.Status,
	}, nil
}
