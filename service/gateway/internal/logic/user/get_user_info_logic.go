package user

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	grpcmw "github.com/luyb177/meow-nook/common/middleware/grpc"
	httpmw "github.com/luyb177/meow-nook/common/middleware/http"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user/v1"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserInfoLogic) GetUserInfo() (*types.UserInfoResp, error) {
	claims, ok := httpmw.ClaimsFromContext(l.ctx)
	if !ok || claims.UserID <= 0 {
		return nil, errorx.ErrUnauthorized
	}

	ctx := grpcmw.InjectRequestID(l.ctx)
	resp, err := l.svcCtx.UserRPC.GetUserInfo(ctx, &userpb.GetUserInfoReq{UserId: claims.UserID})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}
	if resp.User == nil {
		return nil, errorx.ErrUserNotFound
	}

	return &types.UserInfoResp{
		Id:           resp.User.Id,
		Username:     resp.User.Username,
		Avatar:       resp.User.Avatar,
		Phone:        resp.User.Phone,
		Area:         resp.User.Area,
		Gender:       resp.User.Gender,
		Points:       resp.User.Points,
		Role:         resp.User.Role,
		CreatedAt:    resp.User.CreatedAt,
		ServiceTypes: resp.User.ServiceTypes,
	}, nil
}

