package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationSettingsLogic {
	return &GetNotificationSettingsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Notifications
func (l *GetNotificationSettingsLogic) GetNotificationSettings(in *v1.GetNotificationSettingsReq) (*v1.GetNotificationSettingsResp, error) {
	// todo: add your logic here and delete this line

	return &v1.GetNotificationSettingsResp{}, nil
}
