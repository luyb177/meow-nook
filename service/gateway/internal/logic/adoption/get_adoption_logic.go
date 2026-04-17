// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdoptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdoptionLogic {
	return &GetAdoptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdoptionLogic) GetAdoption(req *types.AdoptionPathIdReq) (resp *types.AdoptionApplication, err error) {
	// todo: add your logic here and delete this line

	return
}
