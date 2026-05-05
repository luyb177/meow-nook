package logic

import (
	"context"
	"errors"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetTaskLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogsLogic {
	return &GetTaskLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskLogsLogic) GetTaskLogs(in *v1.GetTaskLogsRequest) (*v1.GetTaskLogsResponse, error) {
	if in.TaskId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "task_id is required", errors.New("task_id is required"))
	}

	logs, err := l.svcCtx.Repo.Task.ListTaskLogs(l.ctx, in.TaskId)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询任务日志失败", err)
	}

	items := make([]*v1.TaskLogVO, 0, len(logs))
	for _, v := range logs {
		item := &v1.TaskLogVO{
			Id:      v.ID,
			TaskId:  v.TaskID,
			UserId:  v.UserID,
			Action:  v.Action,
			Content: v.Content,
		}
		if !v.CreatedAt.IsZero() {
			item.CreatedAt = timestamppb.New(v.CreatedAt)
		}
		items = append(items, item)
	}

	return &v1.GetTaskLogsResponse{
		Items: items,
	}, nil
}
