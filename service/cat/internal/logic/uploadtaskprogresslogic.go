package logic

import (
	"context"
	"encoding/json"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
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
	if _, err := l.svcCtx.Repo.Cat.GetCatTaskByID(l.ctx, in.TaskId); err != nil {
		return nil, errorx.ErrTaskNotFound
	}

	imageURLsJSON, err := json.Marshal(in.ImageUrls)
	if err != nil {
		logger.Error("UploadTaskProgress: marshal image_urls failed")
		return nil, errorx.WrapInternal("序列化图片URL失败", err)
	}

	progress := &cat.CatTaskProgress{
		TaskID:    in.TaskId,
		UserID:    in.UserId,
		Content:   in.Content,
		ImageURLs: string(imageURLsJSON),
	}

	if err := l.svcCtx.Repo.Cat.CreateCatTaskProgress(l.ctx, progress); err != nil {
		logger.Error("UploadTaskProgress: create progress failed")
		return nil, errorx.WrapInternal("上传任务进度失败", err)
	}

	return &v1.UploadTaskProgressResponse{Success: true}, nil
}
