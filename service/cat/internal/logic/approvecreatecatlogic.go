package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveCreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveCreateCatLogic {
	return &ApproveCreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApproveCreateCatLogic) ApproveCreateCat(in *v1.ApproveCreateCatRequest) (*v1.ApproveCreateCatResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ApproveCreateCatResponse{}, nil
}
