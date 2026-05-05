package auth

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	grpcmw "github.com/luyb177/meow-nook/common/middleware/grpc"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user/v1"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (*types.RegisterResp, error) {
	if req.Email == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "邮箱不能为空", errorx.ErrBadRequest)
	}
	if req.Password == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "密码不能为空", errorx.ErrBadRequest)
	}

	ctx := grpcmw.InjectRequestID(l.ctx)
	resp, err := l.svcCtx.UserRPC.Register(ctx, &userpb.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.RegisterResp{UserId: resp.UserId, Token: resp.Token}, nil
}
