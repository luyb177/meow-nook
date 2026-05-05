package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelApplyCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelApplyCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelApplyCreateCatLogic {
	return &CancelApplyCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelApplyCreateCatLogic) CancelApplyCreateCat(in *v1.CancelApplyCreateCatRequest) (*v1.Response, error) {
	err := l.svcCtx.Repo.Cat.CancelApply(
		l.ctx,
		in.ApplyId,
		in.ApplicantUserId,
		in.Reason,
	)

	if err != nil {
		l.Errorf("取消申请失败: %v", err)
		return nil, errorx.WrapDBDelete("取消申请失败", err)
	}

	l.Infof("志愿者取消申请成功, apply_id=%d", in.ApplyId)

	return &v1.Response{}, nil
}
