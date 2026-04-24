package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCatDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCatDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCatDetailLogic {
	return &GetCatDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 猫咪详情
func (l *GetCatDetailLogic) GetCatDetail(in *v1.GetCatDetailRequest) (*v1.GetCatDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.GetCatDetailResponse{}, nil
}
