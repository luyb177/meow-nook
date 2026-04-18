// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	grpcmw "github.com/luyb177/meow-nook/common/middleware/grpc"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user/v1"
)

type TestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestLogic {
	return &TestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestLogic) Test(req *types.TestReq) (resp *types.TestResp, err error) {
	// 注入一下 requestID
	ctx := grpcmw.InjectRequestID(l.ctx)

	res, err := l.svcCtx.UserRPC.Test(ctx, &userpb.TestReq{})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}
	_ = res
	return &types.TestResp{}, nil
}
