package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Profile
func (l *GetUserInfoLogic) GetUserInfo(in *v1.GetUserInfoReq) (*v1.GetUserInfoResp, error) {
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

	return &v1.GetUserInfoResp{
		User: &v1.UserInfo{
			Id:           u.ID,
			Username:     u.Username,
			Avatar:       u.Avatar,
			Phone:        u.Phone,
			Area:         u.Area,
			Gender:       u.Gender,
			Points:       u.Points,
			Role:         u.Role,
			CreatedAt:    u.CreatedAt.Unix(),
			ServiceTypes: u.ServiceTypes,
		},
	}, nil
}
