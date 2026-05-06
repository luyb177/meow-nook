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

type EscalateTaskUrgencyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEscalateTaskUrgencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EscalateTaskUrgencyLogic {
	return &EscalateTaskUrgencyLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *EscalateTaskUrgencyLogic) EscalateTaskUrgency(req *types.EscalateTaskUrgencyReq) (*types.EscalateTaskUrgencyResp, error) {
	logger.Info("EscalateTaskUrgencyLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.EscalateTaskUrgency(l.ctx, &taskpb.EscalateTaskUrgencyRequest{
		TaskId:     req.TaskId,
		OperatorId: uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.EscalateTaskUrgencyResp{
		NewUrgencyLevel: resp.NewUrgencyLevel,
		EscalationCount: resp.EscalationCount,
		Message:         resp.Message,
	}, nil
}
