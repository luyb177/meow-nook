package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/adoption/internal/svc"
	"github.com/luyb177/meow-nook/service/adoption/pb/adoption/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplicationLogic {
	return &GetApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetApplicationLogic) GetApplication(in *v1.GetApplicationReq) (*v1.GetApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &v1.GetApplicationResp{}, nil
}
