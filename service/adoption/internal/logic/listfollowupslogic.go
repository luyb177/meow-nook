package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/adoption/internal/svc"
	"github.com/luyb177/meow-nook/service/adoption/pb/adoption/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFollowUpsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFollowUpsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFollowUpsLogic {
	return &ListFollowUpsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListFollowUpsLogic) ListFollowUps(in *v1.ListFollowUpsReq) (*v1.ListFollowUpsResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ListFollowUpsResp{}, nil
}
