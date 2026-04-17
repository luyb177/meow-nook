package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/post/internal/svc"
	"github.com/luyb177/meow-nook/service/post/pb/post/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoritePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoritePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoritePostLogic {
	return &FavoritePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoritePostLogic) FavoritePost(in *v1.FavoritePostReq) (*v1.FavoritePostResp, error) {
	// todo: add your logic here and delete this line

	return &v1.FavoritePostResp{}, nil
}
