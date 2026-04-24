package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadTaskProgressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadTaskProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadTaskProgressLogic {
	return &UploadTaskProgressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 上传任务进度（图片/备注）
func (l *UploadTaskProgressLogic) UploadTaskProgress(in *v1.UploadTaskProgressRequest) (*v1.UploadTaskProgressResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.UploadTaskProgressResponse{}, nil
}
