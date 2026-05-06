package task

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	taskpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type CancelTaskApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelTaskApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelTaskApplyLogic {
	return &CancelTaskApplyLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CancelTaskApplyLogic) CancelTaskApply(req *types.CancelTaskApplyReq) (*types.CancelTaskApplyResp, error) {
	logger.Info("CancelTaskApplyLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.CancelTaskApply(l.ctx, &taskpb.CancelTaskApplyRequest{
		ApplyId: req.ApplyId,
		Reason:  req.Reason,
		UserId:  uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.CancelTaskApplyResp{
		Message: resp.Message,
	}, nil
}
