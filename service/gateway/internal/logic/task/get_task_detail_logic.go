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

type GetTaskDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskDetailLogic {
	return &GetTaskDetailLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetTaskDetailLogic) GetTaskDetail(req *types.GetTaskDetailReq) (*types.GetTaskDetailResp, error) {
	logger.Info("GetTaskDetailLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.GetTaskDetail(l.ctx, &taskpb.GetTaskDetailRequest{
		TaskId:      req.TaskId,
		RequesterId: uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	out := &types.GetTaskDetailResp{}

	if resp.Task != nil {
		t := resp.Task
		out.Task = types.TaskVO{
			Id:                t.Id,
			CatId:             t.CatId,
			CatName:           t.CatName,
			CatAvatar:         t.CatAvatar,
			ApplyId:           t.ApplyId,
			CreatorId:         t.CreatorId,
			CreatorName:       t.CreatorName,
			Title:             t.Title,
			TaskType:          t.TaskType,
			Summary:           t.Summary,
			Detail:            t.Detail,
			Location:          t.Location,
			Longitude:         t.Longitude,
			Latitude:          t.Latitude,
			Area:              t.Area,
			UrgencyLevel:      t.UrgencyLevel,
			DifficultyLevel:   t.DifficultyLevel,
			RewardPoints:      t.RewardPoints,
			FinalRewardPoints: t.FinalRewardPoints,
			MaxClaimers:       t.MaxClaimers,
			CurrentClaimers:   t.CurrentClaimers,
			Status:            t.Status,
			Remark:            t.Remark,
			DeadlineAt:        logic.PBTimeToString(t.DeadlineAt),
			LastEscalatedAt:   logic.PBTimeToString(t.LastEscalatedAt),
			EscalationCount:   t.EscalationCount,
			CreatedAt:         logic.PBTimeToString(t.CreatedAt),
			UpdatedAt:         logic.PBTimeToString(t.UpdatedAt),
			ImageUrls:         t.ImageUrls,
		}

		tags := make([]types.TagVO, 0, len(t.Tags))
		for _, tg := range t.Tags {
			tags = append(tags, types.TagVO{
				Id:    tg.Id,
				Name:  tg.Name,
				Type:  tg.Type,
				Theme: tg.Theme,
			})
		}
		out.Task.Tags = tags
	}

	if resp.Cat != nil {
		out.Cat = types.CatBriefVO{
			Id:     resp.Cat.Id,
			Name:   resp.Cat.Name,
			Avatar: resp.Cat.Avatar,
			Gender: resp.Cat.Gender,
		}
	}

	return out, nil
}
