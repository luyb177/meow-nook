package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCreateCatApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCreateCatApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCreateCatApplyLogic {
	return &ListCreateCatApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCreateCatApplyLogic) ListCreateCatApply(in *v1.ListCreateCatApplyRequest) (*v1.ListCreateCatApplyResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ListCreateCatApplyResponse{}, nil
}
