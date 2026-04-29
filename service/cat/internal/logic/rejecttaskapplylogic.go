package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type RejectTaskApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectTaskApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectTaskApplyLogic {
	return &RejectTaskApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RejectTaskApplyLogic) RejectTaskApply(in *v1.RejectTaskApplyRequest) (*v1.RejectTaskApplyResponse, error) {
	logger.Info("RejectTaskApply called")

	if err := l.svcCtx.Repo.Task.RejectTaskApply(l.ctx, in.ApplyId, in.ReviewerId, in.RejectReason); err != nil {
		return nil, errorx.WrapDBUpdate("拒绝任务申请失败", err)
	}

	return &v1.RejectTaskApplyResponse{Message: "已拒绝"}, nil
}
