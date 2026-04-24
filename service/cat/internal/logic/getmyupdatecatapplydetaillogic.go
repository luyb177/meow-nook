package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyUpdateCatApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyUpdateCatApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyUpdateCatApplyDetailLogic {
	return &GetMyUpdateCatApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看自己提交的修改申请详情
func (l *GetMyUpdateCatApplyDetailLogic) GetMyUpdateCatApplyDetail(in *v1.GetMyUpdateCatApplyDetailRequest) (*v1.GetMyUpdateCatApplyDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.GetMyUpdateCatApplyDetailResponse{}, nil
}
