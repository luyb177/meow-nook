package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type CancelAdoptApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelAdoptApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelAdoptApplyLogic {
	return &CancelAdoptApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelAdoptApplyLogic) CancelAdoptApply(req *types.CancelAdoptApplyReq) (*types.CancelAdoptApplyResp, error) {
	logger.Info("CancelAdoptApplyLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}
	// todo get userID from token
	userID := uint64(1)

	if req.ApplyId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id is required", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.CancelAdoptApply(l.ctx, &catpb.CancelAdoptApplyRequest{
		ApplyId:     req.ApplyId,
		Reason:      req.Reason,
		ApplicantId: userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.CancelAdoptApplyResp{
		Message: resp.Message,
	}, nil
}
