package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type CancelTaskApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelTaskApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelTaskApplyLogic {
	return &CancelTaskApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelTaskApplyLogic) CancelTaskApply(in *v1.CancelTaskApplyRequest) (*v1.CancelTaskApplyResponse, error) {
	logger.Info("CancelTaskApply called")

	if err := l.svcCtx.Repo.Task.CancelTaskApply(l.ctx, in.ApplyId, in.UserId, in.Reason); err != nil {
		return nil, errorx.WrapDBUpdate("取消任务申请失败", err)
	}

	return &v1.CancelTaskApplyResponse{Message: "已取消申请"}, nil
}
