package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *v1.ChangePasswordReq) (*v1.ChangePasswordResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ChangePasswordResp{}, nil
}
