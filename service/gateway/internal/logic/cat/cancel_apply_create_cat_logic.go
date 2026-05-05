// internal/logic/cat/cancel_apply_create_cat_logic.go
package cat

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
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
	//uid := getUserID(l.ctx)
	//if uid == 0 {
	//	return nil, errorx.ErrUnauthorized
	//}
	// todo 从 JWT 中获取用户 ID
	uid := uint64(1)

	_, err := l.svcCtx.CatRPC.CancelApplyCreateCat(l.ctx, &catpb.CancelApplyCreateCatRequest{
		ApplyId:         req.Id,
		ApplicantUserId: uid,
		Reason:          req.Reason,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.Response{}, nil
}
