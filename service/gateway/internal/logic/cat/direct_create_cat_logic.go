// internal/logic/cat/direct_create_cat_logic.go
package cat

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type DirectCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDirectCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DirectCreateCatLogic {
	return &DirectCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DirectCreateCatLogic) DirectCreateCat(req *types.DirectCreateCatReq) (*types.DirectCreateCatResp, error) {
	//if !isAdmin(l.ctx) {
	//	return nil, errorx.ErrPermissionDenied
	//}
	//
	//uid := getUserID(l.ctx)
	//if uid == 0 {
	//	return nil, errorx.ErrUnauthorized
	//}

	// todo  从JWT中获取用户ID
	uid := uint64(1)

	// todo casbin 权限检测

	resp, err := l.svcCtx.CatRPC.DirectCreateCat(l.ctx, &catpb.DirectCreateCatRequest{
		CatCode:                 req.CatCode,
		Name:                    req.Name,
		Breed:                   req.Breed,
		Color:                   req.Color,
		Gender:                  req.Gender,
		BodySize:                req.BodySize,
		AgeStage:                req.AgeStage,
		Description:             req.Description,
		DiscoveryAddress:        req.DiscoveryAddress,
		Longitude:               req.Longitude,
		Latitude:                req.Latitude,
		IsVaccinated:            req.IsVaccinated,
		IsHealthy:               req.IsHealthy,
		NeedMedicalIntervention: req.NeedMedicalIntervention,
		SterilizationStatus:     req.SterilizationStatus,
		Images:                  convertImagesToPB(req.Images),
		TagIds:                  req.TagIds,
		OperatorId:              uid,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.DirectCreateCatResp{
		CatId:   resp.CatId,
		CatCode: resp.CatCode,
	}, nil
}
