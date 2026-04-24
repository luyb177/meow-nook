package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCreateCatLogic {
	return &ApplyCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请创建小猫档案
func (l *ApplyCreateCatLogic) ApplyCreateCat(in *v1.ApplyCreateCatRequest) (*v1.ApplyCreateCatResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ApplyCreateCatResponse{}, nil
}
