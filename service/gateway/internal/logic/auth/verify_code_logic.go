package auth

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	grpcmw "github.com/luyb177/meow-nook/common/middleware/grpc"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user/v1"
)

type VerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCodeLogic {
	return &VerifyCodeLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *VerifyCodeLogic) VerifyCode(req *types.VerifyCodeReq) (*types.Response, error) {
	if req.Target == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "目标地址不能为空", errorx.ErrBadRequest)
	}
	if req.Code == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "验证码不能为空", errorx.ErrBadRequest)
	}

	ctx := grpcmw.InjectRequestID(l.ctx)
	_, err := l.svcCtx.UserRPC.VerifyCode(ctx, &userpb.VerifyCodeReq{
		Target:  req.Target,
		Channel: userpb.VerifyChannel(req.Channel),
		Purpose: userpb.VerifyPurpose(req.Purpose),
		Code:    req.Code,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.Response{}, nil
}
