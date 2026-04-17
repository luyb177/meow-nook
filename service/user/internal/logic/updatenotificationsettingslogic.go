package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNotificationSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotificationSettingsLogic {
	return &UpdateNotificationSettingsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateNotificationSettingsLogic) UpdateNotificationSettings(in *v1.UpdateNotificationSettingsReq) (*v1.UpdateNotificationSettingsResp, error) {
	// todo: add your logic here and delete this line

	return &v1.UpdateNotificationSettingsResp{}, nil
}
