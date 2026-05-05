package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"go.uber.org/zap"
)

type RejectCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectCreateCatLogic {
	return &RejectCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RejectCreateCatLogic) RejectCreateCat(in *v1.RejectCreateCatRequest) (*v1.Response, error) {
	err := l.svcCtx.Repo.Cat.RejectApply(
		l.ctx,
		in.ApplyId,
		in.OperatorId,
		in.Reason,
	)

	if err != nil {
		logger.Error("驳回申请失败", zap.Error(err))
		return nil, errorx.WrapDBUpdate("驳回申请失败", err)
	}

	logger.Info("管理员驳回申请成功", zap.Uint64("apply_id", in.ApplyId), zap.String("reason", in.Reason))

	return &v1.Response{}, nil
}
