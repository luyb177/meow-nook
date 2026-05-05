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

type ApproveTaskApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveTaskApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveTaskApplyLogic {
	return &ApproveTaskApplyLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ApproveTaskApplyLogic) ApproveTaskApply(req *types.ApproveTaskApplyReq) (*types.ApproveTaskApplyResp, error) {
	logger.Info("ApproveTaskApplyLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.ApproveTaskApply(l.ctx, &taskpb.ApproveTaskApplyRequest{
		ApplyId:         req.ApplyId,
		UrgencyLevel:    req.UrgencyLevel,
		DifficultyLevel: req.DifficultyLevel,
		RewardPoints:    req.RewardPoints,
		MaxClaimers:     req.MaxClaimers,
		Area:            req.Area,
		Remark:          req.Remark,
		ReviewerId:      userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.ApproveTaskApplyResp{
		TaskId:  resp.TaskId,
		Message: resp.Message,
	}, nil
}
