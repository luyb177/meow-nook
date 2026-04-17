package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/post/internal/svc"
	"github.com/luyb177/meow-nook/service/post/pb/post/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCommentLogic) AddComment(in *v1.AddCommentReq) (*v1.AddCommentResp, error) {
	// todo: add your logic here and delete this line

	return &v1.AddCommentResp{}, nil
}
