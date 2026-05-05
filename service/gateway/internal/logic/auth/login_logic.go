package auth

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	grpcmw "github.com/luyb177/meow-nook/common/middleware/grpc"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user/v1"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	if req.Email == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "邮箱不能为空", errorx.ErrBadRequest)
	}
	if req.Password == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "密码不能为空", errorx.ErrBadRequest)
	}

	ctx := grpcmw.InjectRequestID(l.ctx)
	resp, err := l.svcCtx.UserRPC.Login(ctx, &userpb.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.LoginResp{
		Token:  resp.Token,
		UserId: resp.UserId,
	}, nil
}
