package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/post/internal/svc"
	"github.com/luyb177/meow-nook/service/post/pb/post/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostLogic) GetPost(in *v1.GetPostReq) (*v1.GetPostResp, error) {
	// todo: add your logic here and delete this line

	return &v1.GetPostResp{}, nil
}
