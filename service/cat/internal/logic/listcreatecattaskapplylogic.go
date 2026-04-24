package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCreateCatTaskApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCreateCatTaskApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCreateCatTaskApplyLogic {
	return &ListCreateCatTaskApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCreateCatTaskApplyLogic) ListCreateCatTaskApply(in *v1.ListCreateCatTaskApplyRequest) (*v1.ListCreateCatTaskApplyResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ListCreateCatTaskApplyResponse{}, nil
}
