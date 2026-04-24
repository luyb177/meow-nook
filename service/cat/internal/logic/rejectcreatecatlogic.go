package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectCreateCatLogic {
	return &RejectCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectCreateCatLogic) RejectCreateCat(in *v1.RejectCreateCatRequest) (*v1.RejectCreateCatResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.RejectCreateCatResponse{}, nil
}
