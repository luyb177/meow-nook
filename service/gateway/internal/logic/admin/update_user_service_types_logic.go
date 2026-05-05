package admin

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	grpcmw "github.com/luyb177/meow-nook/common/middleware/grpc"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user/v1"
)

type UpdateUserServiceTypesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserServiceTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserServiceTypesLogic {
	return &UpdateUserServiceTypesLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserServiceTypesLogic) UpdateUserServiceTypes(req *types.AdminUpdateUserServiceTypesReq) (*types.Response, error) {
	if req.UserId <= 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "userId 不能为空", errorx.ErrBadRequest)
	}

	ctx := grpcmw.InjectRequestID(l.ctx)
	if _, err := l.svcCtx.UserRPC.AdminUpdateServiceTypes(ctx, &userpb.AdminUpdateServiceTypesReq{
		UserId:        req.UserId,
		ServiceTypes:  req.ServiceTypes,
	}); err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.Response{}, nil
}

