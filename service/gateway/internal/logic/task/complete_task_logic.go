package task

import (
	"context"
	"encoding/json"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	taskpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type CompleteTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteTaskLogic {
	return &CompleteTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CompleteTaskLogic) CompleteTask(req *types.CompleteTaskReq) (*types.CompleteTaskResp, error) {
	logger.Info("CompleteTaskLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	// 图片URL列表转JSON
	imageURLsJSON, _ := json.Marshal(req.ImageUrls)

	resp, err := l.svcCtx.CatRPC.CompleteTask(l.ctx, &taskpb.CompleteTaskRequest{
		ClaimId:   req.ClaimId,
		Content:   req.Content,
		ImageUrls: []string{string(imageURLsJSON)},
		UserId:    uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.CompleteTaskResp{
		RewardPoints:        resp.RewardPoints,
		IsOverdue:           resp.IsOverdue,
		ReputationDeduction: resp.ReputationDeduction,
		Message:             resp.Message,
	}, nil
}
