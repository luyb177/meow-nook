package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/adoption/internal/svc"
	"github.com/luyb177/meow-nook/service/adoption/pb/adoption/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyAdoptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAdoptionLogic {
	return &ApplyAdoptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyAdoptionLogic) ApplyAdoption(in *v1.ApplyAdoptionReq) (*v1.ApplyAdoptionResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ApplyAdoptionResp{}, nil
}
