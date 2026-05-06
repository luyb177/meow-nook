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

type GetTaskFlowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskFlowsLogic {
	return &GetTaskFlowsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetTaskFlowsLogic) GetTaskFlows(req *types.GetTaskFlowsReq) (*types.GetTaskFlowsResp, error) {
	logger.Info("GetTaskFlowsLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.GetTaskFlows(l.ctx, &taskpb.GetTaskFlowsRequest{
		TaskId:      req.TaskId,
		RequesterId: uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.TaskFlowVO, 0, len(resp.Items))
	for _, f := range resp.Items {
		items = append(items, types.TaskFlowVO{
			Id:         f.Id,
			TaskId:     f.TaskId,
			UserId:     f.UserId,
			UserName:   f.UserName,
			Action:     f.Action,
			FromStatus: f.FromStatus,
			ToStatus:   f.ToStatus,
			Remark:     f.Remark,
			CreatedAt:  logic.PBTimeToString(f.CreatedAt),
		})
	}

	return &types.GetTaskFlowsResp{
		Items: items,
	}, nil
}
