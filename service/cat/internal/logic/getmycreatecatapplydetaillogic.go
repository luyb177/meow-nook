package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyCreateCatApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyCreateCatApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyCreateCatApplyDetailLogic {
	return &GetMyCreateCatApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看自己提交的创建申请详情
func (l *GetMyCreateCatApplyDetailLogic) GetMyCreateCatApplyDetail(in *v1.GetMyCreateCatApplyDetailRequest) (*v1.GetMyCreateCatApplyDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.GetMyCreateCatApplyDetailResponse{}, nil
}
