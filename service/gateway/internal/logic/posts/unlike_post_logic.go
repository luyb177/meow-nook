// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package posts

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnlikePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnlikePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikePostLogic {
	return &UnlikePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnlikePostLogic) UnlikePost(req *types.PostPathIdReq) error {
	// todo: add your logic here and delete this line

	return nil
}
