// internal/logic/cat/cancel_apply_create_cat_logic.go
package cat

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type CancelApplyCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelApplyCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelApplyCreateCatLogic {
	return &CancelApplyCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelApplyCreateCatLogic) CancelApplyCreateCat(req *types.CancelApplyCreateCatReq) (*types.Response, error) {
	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	_, err = l.svcCtx.CatRPC.CancelApplyCreateCat(l.ctx, &catpb.CancelApplyCreateCatRequest{
		ApplyId:         req.Id,
		ApplicantUserId: uint64(userID),
		Reason:          req.Reason,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.Response{}, nil
}
