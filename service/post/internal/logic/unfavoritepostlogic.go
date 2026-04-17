package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/post/internal/svc"
	"github.com/luyb177/meow-nook/service/post/pb/post/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfavoritePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnfavoritePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfavoritePostLogic {
	return &UnfavoritePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnfavoritePostLogic) UnfavoritePost(in *v1.UnfavoritePostReq) (*v1.UnfavoritePostResp, error) {
	// todo: add your logic here and delete this line

	return &v1.UnfavoritePostResp{}, nil
}
