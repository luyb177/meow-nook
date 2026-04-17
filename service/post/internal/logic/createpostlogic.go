package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/post/internal/svc"
	"github.com/luyb177/meow-nook/service/post/pb/post/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(in *v1.CreatePostReq) (*v1.CreatePostResp, error) {
	// todo: add your logic here and delete this line

	return &v1.CreatePostResp{}, nil
}
