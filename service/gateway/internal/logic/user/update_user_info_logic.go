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

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (*types.Response, error) {
	claims, ok := httpmw.ClaimsFromContext(l.ctx)
	if !ok || claims.UserID <= 0 {
		return nil, errorx.ErrUnauthorized
	}

	pbReq := &userpb.UpdateUserInfoReq{UserId: claims.UserID}
	if req.Username != nil {
		pbReq.Username = *req.Username
	}
	if req.Avatar != nil {
		pbReq.Avatar = *req.Avatar
	}
	if req.Phone != nil {
		pbReq.Phone = *req.Phone
	}
	if req.Area != nil {
		pbReq.Area = *req.Area
	}

	ctx := grpcmw.InjectRequestID(l.ctx)
	if _, err := l.svcCtx.UserRPC.UpdateUserInfo(ctx, pbReq); err != nil {
		return nil, errorx.FromGRPC(err)
	}
	return &types.Response{}, nil
}

