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

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateTaskLogic) UpdateTask(req *types.UpdateTaskReq) (*types.UpdateTaskResp, error) {
	logger.Info("UpdateTaskLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.UpdateTask(l.ctx, &taskpb.UpdateTaskRequest{
		TaskId:          req.TaskId,
		Title:           req.Title,
		Summary:         req.Summary,
		Detail:          req.Detail,
		DifficultyLevel: req.DifficultyLevel,
		RewardPoints:    req.RewardPoints,
		DeadlineAt:      logic.ParseTimePtr(req.DeadlineAt),
		Remark:          req.Remark,
		OperatorId:      uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.UpdateTaskResp{
		Message: resp.Message,
	}, nil
}
