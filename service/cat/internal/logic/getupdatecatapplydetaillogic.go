package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUpdateCatApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUpdateCatApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUpdateCatApplyDetailLogic {
	return &GetUpdateCatApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUpdateCatApplyDetailLogic) GetUpdateCatApplyDetail(in *v1.GetUpdateCatApplyDetailRequest) (*v1.GetUpdateCatApplyDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.GetUpdateCatApplyDetailResponse{}, nil
}
