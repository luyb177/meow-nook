// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package cats

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHeatmapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHeatmapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHeatmapLogic {
	return &GetHeatmapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHeatmapLogic) GetHeatmap(req *types.HeatmapReq) (resp *types.HeatmapResp, err error) {
	// todo: add your logic here and delete this line

	return
}
