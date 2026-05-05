package logic

import (
	"context"
	"errors"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeductCreditForTaskViolationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductCreditForTaskViolationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductCreditForTaskViolationLogic {
	return &DeductCreditForTaskViolationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductCreditForTaskViolationLogic) DeductCreditForTaskViolation(in *v1.DeductCreditForTaskViolationReq) (*v1.DeductCreditForTaskViolationResp, error) {
	if in.UserId <= 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "user_id 不能为空", errorx.ErrBadRequest)
	}

	newPoints, err := l.svcCtx.Repo.User.AddPointsDelta(l.ctx, in.UserId, -50)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	if err != nil {
		return nil, errorx.WrapInternal("扣除信誉积分失败", err)
	}

	return &v1.DeductCreditForTaskViolationResp{CreditPoints: newPoints}, nil
}
