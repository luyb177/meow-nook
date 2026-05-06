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

type RejectTaskApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectTaskApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectTaskApplyLogic {
	return &RejectTaskApplyLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *RejectTaskApplyLogic) RejectTaskApply(req *types.RejectTaskApplyReq) (*types.RejectTaskApplyResp, error) {
	logger.Info("RejectTaskApplyLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.RejectTaskApply(l.ctx, &taskpb.RejectTaskApplyRequest{
		ApplyId:      req.ApplyId,
		RejectReason: req.RejectReason,
		ReviewerId:   uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.RejectTaskApplyResp{
		Message: resp.Message,
	}, nil
}
