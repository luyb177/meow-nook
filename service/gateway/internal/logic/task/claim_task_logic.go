package task

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	taskpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ClaimTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClaimTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimTaskLogic {
	return &ClaimTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ClaimTaskLogic) ClaimTask(req *types.ClaimTaskReq) (*types.ClaimTaskResp, error) {
	logger.Info("ClaimTaskLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo g u f t
	userID := uint64(1)

	resp, err := l.svcCtx.CatRPC.ClaimTask(l.ctx, &taskpb.ClaimTaskRequest{
		TaskId: req.TaskId,
		UserId: userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.ClaimTaskResp{
		ClaimId: resp.ClaimId,
		Message: resp.Message,
	}, nil
}
