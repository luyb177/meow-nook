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

type ApplyCreateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCreateTaskLogic {
	return &ApplyCreateTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ApplyCreateTaskLogic) ApplyCreateTask(req *types.ApplyCreateTaskReq) (*types.ApplyCreateTaskResp, error) {
	logger.Info("ApplyCreateTaskLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.ApplyCreateTask(l.ctx, &taskpb.ApplyCreateTaskRequest{
		CatId:           req.CatId,
		Title:           req.Title,
		TaskType:        req.TaskType,
		Summary:         req.Summary,
		Detail:          req.Detail,
		Location:        req.Location,
		Longitude:       req.Longitude,
		Latitude:        req.Latitude,
		DeadlineAt:      logic.ParseTimePtr(req.DeadlineAt),
		ImageUrls:       req.ImageUrls,
		TagIds:          req.TagIds,
		ApplicantUserId: uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.ApplyCreateTaskResp{
		ApplyId: resp.ApplyId,
		Status:  resp.Status,
		Message: resp.Message,
	}, nil
}
