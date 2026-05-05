package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type ApproveTaskApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveTaskApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveTaskApplyLogic {
	return &ApproveTaskApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveTaskApplyLogic) ApproveTaskApply(in *v1.ApproveTaskApplyRequest) (*v1.ApproveTaskApplyResponse, error) {
	logger.Info("ApproveTaskApply called")

	task, err := l.svcCtx.Repo.Task.ApproveTaskApply(l.ctx, in.ApplyId, in.ReviewerId)
	if err != nil {
		return nil, errorx.WrapDBUpdate("审核任务申请失败", err)
	}

	return &v1.ApproveTaskApplyResponse{
		TaskId:  task.ID,
		Message: "审核通过",
	}, nil
}
