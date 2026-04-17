// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package cats

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCatStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCatStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCatStatsLogic {
	return &GetCatStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCatStatsLogic) GetCatStats() (resp *types.CatStatsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
