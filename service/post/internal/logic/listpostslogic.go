package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/post/internal/svc"
	"github.com/luyb177/meow-nook/service/post/pb/post/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPostsLogic) ListPosts(in *v1.ListPostsReq) (*v1.ListPostsResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ListPostsResp{}, nil
}
