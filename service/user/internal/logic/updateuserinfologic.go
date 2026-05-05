package logic

import (
	"context"
	"unicode/utf8"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *v1.UpdateUserInfoReq) (*v1.UpdateUserInfoResp, error) {
	if in.UserId <= 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "user_id 不能为空", errorx.ErrBadRequest)
	}

	fields := make(map[string]any)

	if in.Username != "" {
		if utf8.RuneCountInString(in.Username) > 64 {
			return nil, errorx.New(errorx.CodeBadRequest, "姓名长度不能超过64位")
		}
		fields["username"] = in.Username
	}
	if in.Avatar != "" {
		if utf8.RuneCountInString(in.Avatar) > 512 {
			return nil, errorx.New(errorx.CodeBadRequest, "头像地址长度不能超过512位")
		}
		fields["avatar"] = in.Avatar
	}
	if in.Phone != "" {
		if utf8.RuneCountInString(in.Phone) > 32 {
			return nil, errorx.New(errorx.CodeBadRequest, "手机号长度不能超过32位")
		}
		fields["phone"] = in.Phone
	}
	if in.Area != "" {
		if utf8.RuneCountInString(in.Area) > 128 {
			return nil, errorx.New(errorx.CodeBadRequest, "地区长度不能超过128位")
		}
		fields["area"] = in.Area
	}
	if in.Gender != "" {
		if utf8.RuneCountInString(in.Gender) > 16 {
			return nil, errorx.New(errorx.CodeBadRequest, "性别长度不能超过16位")
		}
		fields["gender"] = in.Gender
	}

	if len(fields) == 0 {
		return &v1.UpdateUserInfoResp{}, nil
	}

	u, err := l.svcCtx.Repo.User.FindByID(l.ctx, in.UserId)
	if err != nil {
		return nil, errorx.WrapInternal("查询用户失败", err)
	}
	if u == nil {
		return nil, errorx.ErrUserNotFound
	}

	if err := l.svcCtx.Repo.User.UpdateFields(l.ctx, in.UserId, fields); err != nil {
		return nil, errorx.WrapInternal("更新用户信息失败", err)
	}

	return &v1.UpdateUserInfoResp{}, nil
}
