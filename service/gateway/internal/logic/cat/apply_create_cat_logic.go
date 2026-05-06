// internal/logic/cat/apply_create_cat_logic.go
package cat

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
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
	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

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
		ApplicantUserId:  uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.ApplyCreateCatResp{
		ApplyId:   resp.ApplyId,
		CreatedAt: resp.CreatedAt.AsTime().Format(time.DateTime),
	}, nil
}
