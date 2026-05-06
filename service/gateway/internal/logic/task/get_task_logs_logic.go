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

type GetTaskLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogsLogic {
	return &GetTaskLogsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetTaskLogsLogic) GetTaskLogs(req *types.GetTaskLogsReq) (*types.GetTaskLogsResp, error) {
	logger.Info("GetTaskLogsLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.GetTaskLogs(l.ctx, &taskpb.GetTaskLogsRequest{
		TaskId:      req.TaskId,
		RequesterId: uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.TaskLogVO, 0, len(resp.Items))
	for _, log := range resp.Items {
		items = append(items, types.TaskLogVO{
			Id:        log.Id,
			TaskId:    log.TaskId,
			UserId:    log.UserId,
			UserName:  log.UserName,
			Action:    log.Action,
			Content:   log.Content,
			CreatedAt: logic.PBTimeToString(log.CreatedAt),
		})
	}

	return &types.GetTaskLogsResp{
		Items: items,
	}, nil
}
