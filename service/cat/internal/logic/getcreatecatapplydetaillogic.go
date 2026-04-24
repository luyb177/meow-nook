package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCreateCatApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCreateCatApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCreateCatApplyDetailLogic {
	return &GetCreateCatApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCreateCatApplyDetailLogic) GetCreateCatApplyDetail(in *v1.GetCreateCatApplyDetailRequest) (*v1.GetCreateCatApplyDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.GetCreateCatApplyDetailResponse{}, nil
}
