package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCreateCatTaskApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCreateCatTaskApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCreateCatTaskApplyDetailLogic {
	return &GetCreateCatTaskApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCreateCatTaskApplyDetailLogic) GetCreateCatTaskApplyDetail(in *v1.GetCreateCatTaskApplyDetailRequest) (*v1.GetCreateCatTaskApplyDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.GetCreateCatTaskApplyDetailResponse{}, nil
}
