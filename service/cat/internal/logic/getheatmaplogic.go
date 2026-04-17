package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHeatmapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHeatmapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHeatmapLogic {
	return &GetHeatmapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetHeatmapLogic) GetHeatmap(in *v1.GetHeatmapReq) (*v1.GetHeatmapResp, error) {
	// todo: add your logic here and delete this line

	return &v1.GetHeatmapResp{}, nil
}
