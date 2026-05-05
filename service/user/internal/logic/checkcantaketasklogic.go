package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckCanTakeTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckCanTakeTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckCanTakeTaskLogic {
	return &CheckCanTakeTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckCanTakeTaskLogic) CheckCanTakeTask(in *v1.CheckCanTakeTaskReq) (*v1.CheckCanTakeTaskResp, error) {
	if in.UserId <= 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "user_id 不能为空", errorx.ErrBadRequest)
	}

	u, err := l.svcCtx.Repo.User.FindByID(l.ctx, in.UserId)
	if err != nil {
		return nil, errorx.WrapInternal("查询用户失败", err)
	}
	if u == nil {
		return nil, errorx.ErrUserNotFound
	}
	if u.Points < 50 {
		return nil, errorx.New(errorx.CodeForbidden, "信誉积分低于50，无法接任务")
	}

	return &v1.CheckCanTakeTaskResp{CreditPoints: u.Points}, nil
}
