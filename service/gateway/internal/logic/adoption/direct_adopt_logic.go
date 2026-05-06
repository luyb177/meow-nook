// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DirectAdoptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDirectAdoptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DirectAdoptLogic {
	return &DirectAdoptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DirectAdopt
func (l *DirectAdoptLogic) DirectAdopt(req *types.DirectAdoptReq) (*types.DirectAdoptResp, error) {
	logger.Info("DirectAdoptLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	if req.CatId == 0 || req.AdopterId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "cat_id 和 adopter_id 必填", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.DirectAdopt(l.ctx, &catpb.DirectAdoptRequest{
		CatId:       req.CatId,
		AdopterId:   req.AdopterId,
		AgreementNo: req.AgreementNo,
		Note:        req.Note,
		CreatorId:   uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.DirectAdoptResp{
		AdoptionId: resp.AdoptionId,
		Message:    resp.Message,
	}, nil
}
