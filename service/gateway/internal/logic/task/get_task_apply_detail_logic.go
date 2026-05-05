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

type GetTaskApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskApplyDetailLogic {
	return &GetTaskApplyDetailLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetTaskApplyDetailLogic) GetTaskApplyDetail(req *types.GetTaskApplyDetailReq) (*types.GetTaskApplyDetailResp, error) {
	logger.Info("GetTaskApplyDetailLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo g u f t
	userID := uint64(1)

	resp, err := l.svcCtx.CatRPC.GetTaskApplyDetail(l.ctx, &taskpb.GetTaskApplyDetailRequest{
		ApplyId:     req.ApplyId,
		RequesterId: userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	out := &types.GetTaskApplyDetailResp{}

	if resp.Apply != nil {
		a := resp.Apply
		out.Apply = types.TaskApplyVO{
			Id:              a.Id,
			CatId:           a.CatId,
			CatName:         a.CatName,
			CatAvatar:       a.CatAvatar,
			ApplicantUserId: a.ApplicantUserId,
			ApplicantName:   a.ApplicantName,
			ApplicantAvatar: a.ApplicantAvatar,
			Title:           a.Title,
			TaskType:        a.TaskType,
			Summary:         a.Summary,
			Detail:          a.Detail,
			Location:        a.Location,
			Longitude:       a.Longitude,
			Latitude:        a.Latitude,
			Status:          a.Status,
			ReviewerId:      a.ReviewerId,
			ReviewerName:    a.ReviewerName,
			ReviewedAt:      logic.PBTimeToString(a.ReviewedAt),
			RejectReason:    a.RejectReason,
			UrgencyLevel:    a.UrgencyLevel,
			DifficultyLevel: a.DifficultyLevel,
			RewardPoints:    a.RewardPoints,
			TaskId:          a.TaskId,
			DeadlineAt:      logic.PBTimeToString(a.DeadlineAt),
			CreatedAt:       logic.PBTimeToString(a.CreatedAt),
			UpdatedAt:       logic.PBTimeToString(a.UpdatedAt),
		}

		// 填充图片
		imageURLs := make([]string, 0, len(a.ImageUrls))
		for _, url := range a.ImageUrls {
			imageURLs = append(imageURLs, url)
		}
		out.Apply.ImageUrls = imageURLs

		// 填充标签
		tags := make([]types.TagVO, 0, len(a.Tags))
		for _, t := range a.Tags {
			tags = append(tags, types.TagVO{
				Id:    t.Id,
				Name:  t.Name,
				Type:  t.Type,
				Theme: t.Theme,
			})
		}
		out.Apply.Tags = tags
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
