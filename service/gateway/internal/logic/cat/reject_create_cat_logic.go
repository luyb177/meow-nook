// internal/logic/cat/reject_create_cat_logic.go
package cat

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type RejectCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectCreateCatLogic {
	return &RejectCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RejectCreateCatLogic) RejectCreateCat(req *types.RejectCreateCatReq) (*types.Response, error) {
	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	// todo casbin 权限检测

	_, err = l.svcCtx.CatRPC.RejectCreateCat(l.ctx, &catpb.RejectCreateCatRequest{
		ApplyId:    req.Id,
		Reason:     req.Reason,
		OperatorId: uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.Response{}, nil
}
