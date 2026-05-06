package task

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	taskpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type DirectCreateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDirectCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DirectCreateTaskLogic {
	return &DirectCreateTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DirectCreateTaskLogic) DirectCreateTask(req *types.DirectCreateTaskReq) (*types.DirectCreateTaskResp, error) {
	logger.Info("DirectCreateTaskLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.DirectCreateTask(l.ctx, &taskpb.DirectCreateTaskRequest{
		CatId:           req.CatId,
		Title:           req.Title,
		TaskType:        req.TaskType,
		Summary:         req.Summary,
		Detail:          req.Detail,
		Location:        req.Location,
		Longitude:       req.Longitude,
		Latitude:        req.Latitude,
		Area:            req.Area,
		UrgencyLevel:    req.UrgencyLevel,
		DifficultyLevel: req.DifficultyLevel,
		RewardPoints:    req.RewardPoints,
		MaxClaimers:     req.MaxClaimers,
		DeadlineAt:      logic.ParseTimePtr(req.DeadlineAt),
		Remark:          req.Remark,
		ImageUrls:       req.ImageUrls,
		TagIds:          req.TagIds,
		CreatorId:       uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.DirectCreateTaskResp{
		TaskId:  resp.TaskId,
		Message: resp.Message,
	}, nil
}
