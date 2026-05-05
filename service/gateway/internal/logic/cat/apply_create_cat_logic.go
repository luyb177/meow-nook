// internal/logic/cat/apply_create_cat_logic.go
package cat

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCreateCatLogic {
	return &ApplyCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyCreateCatLogic) ApplyCreateCat(req *types.ApplyCreateCatReq) (*types.ApplyCreateCatResp, error) {
	// todo 从 JWT 中获取用户 ID
	//uid := getUserID(l.ctx)
	//if uid == 0 {
	//	return nil, errorx.ErrUnauthorized
	//}
	uid := uint64(1)

	resp, err := l.svcCtx.CatRPC.ApplyCreateCat(l.ctx, &catpb.ApplyCreateCatRequest{
		Name:             req.Name,
		Gender:           req.Gender,
		BodySize:         req.BodySize,
		AgeStage:         req.AgeStage,
		Description:      req.Description,
		DiscoveryAddress: req.DiscoveryAddress,
		Longitude:        req.Longitude,
		Latitude:         req.Latitude,
		Images:           convertImagesToPB(req.Images),
		TagIds:           req.TagIds,
		ApplicantUserId:  uid,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.ApplyCreateCatResp{
		ApplyId:   resp.ApplyId,
		CreatedAt: resp.CreatedAt.AsTime().Format(time.DateTime),
	}, nil
}
