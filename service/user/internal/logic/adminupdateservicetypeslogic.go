package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateServiceTypesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminUpdateServiceTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateServiceTypesLogic {
	return &AdminUpdateServiceTypesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Admin
func (l *AdminUpdateServiceTypesLogic) AdminUpdateServiceTypes(in *v1.AdminUpdateServiceTypesReq) (*v1.AdminUpdateServiceTypesResp, error) {
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

	if err := l.svcCtx.Repo.User.UpdateFields(l.ctx, in.UserId, map[string]any{
		"service_types": in.ServiceTypes,
	}); err != nil {
		return nil, errorx.WrapInternal("更新可服务类型失败", err)
	}

	return &v1.AdminUpdateServiceTypesResp{}, nil
}
