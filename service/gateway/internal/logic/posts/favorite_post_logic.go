// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package posts

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoritePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoritePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoritePostLogic {
	return &FavoritePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoritePostLogic) FavoritePost(req *types.PostPathIdReq) error {
	// todo: add your logic here and delete this line

	return nil
}
