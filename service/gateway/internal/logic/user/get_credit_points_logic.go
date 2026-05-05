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

type GetCreditPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCreditPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCreditPointsLogic {
	return &GetCreditPointsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetCreditPointsLogic) GetCreditPoints() (*types.CreditPointsResp, error) {
	claims, ok := httpmw.ClaimsFromContext(l.ctx)
	if !ok || claims.UserID <= 0 {
		return nil, errorx.ErrUnauthorized
	}

	ctx := grpcmw.InjectRequestID(l.ctx)
	resp, err := l.svcCtx.UserRPC.GetCreditPoints(ctx, &userpb.GetCreditPointsReq{UserId: claims.UserID})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}
	return &types.CreditPointsResp{CreditPoints: resp.CreditPoints}, nil
}
