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

type CancelTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelTaskLogic {
	return &CancelTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CancelTaskLogic) CancelTask(req *types.CancelTaskReq) (*types.CancelTaskResp, error) {
	logger.Info("CancelTaskLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.CancelTask(l.ctx, &taskpb.CancelTaskRequest{
		TaskId:     req.TaskId,
		Reason:     req.Reason,
		OperatorId: userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.CancelTaskResp{
		Message: resp.Message,
	}, nil
}
