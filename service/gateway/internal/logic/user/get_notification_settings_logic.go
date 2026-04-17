// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationSettingsLogic {
	return &GetNotificationSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotificationSettingsLogic) GetNotificationSettings() (resp *types.NotificationSettings, err error) {
	// todo: add your logic here and delete this line

	return
}
