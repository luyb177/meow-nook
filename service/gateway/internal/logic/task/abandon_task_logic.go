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

type AbandonTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAbandonTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbandonTaskLogic {
	return &AbandonTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AbandonTaskLogic) AbandonTask(req *types.AbandonTaskReq) (*types.AbandonTaskResp, error) {
	logger.Info("AbandonTaskLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.AbandonTask(l.ctx, &taskpb.AbandonTaskRequest{
		ClaimId: req.ClaimId,
		Reason:  req.Reason,
		UserId:  uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.AbandonTaskResp{
		Message: resp.Message,
	}, nil
}
