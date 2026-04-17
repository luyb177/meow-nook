// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAdoptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAdoptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAdoptionsLogic {
	return &ListAdoptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAdoptionsLogic) ListAdoptions(req *types.ListAdoptionReq) (resp *types.ListAdoptionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
