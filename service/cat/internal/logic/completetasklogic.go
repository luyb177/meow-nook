package logic

import (
	"context"
	"encoding/json"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type CompleteTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteTaskLogic {
	return &CompleteTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteTaskLogic) CompleteTask(in *v1.CompleteTaskRequest) (*v1.CompleteTaskResponse, error) {
	logger.Info("CompleteTask called")

	imgJSON, _ := json.Marshal(in.ImageUrls)

	if err := l.svcCtx.Repo.Task.CompleteTaskClaim(l.ctx, in.ClaimId, in.UserId, in.Content, string(imgJSON)); err != nil {
		return nil, errorx.WrapDBUpdate("完成任务失败", err)
	}

	return &v1.CompleteTaskResponse{Message: "任务完成"}, nil
}
