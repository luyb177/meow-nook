package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/adoption/internal/svc"
	"github.com/luyb177/meow-nook/service/adoption/pb/adoption/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListApplicationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApplicationsLogic {
	return &ListApplicationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListApplicationsLogic) ListApplications(in *v1.ListApplicationsReq) (*v1.ListApplicationsResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ListApplicationsResp{}, nil
}
