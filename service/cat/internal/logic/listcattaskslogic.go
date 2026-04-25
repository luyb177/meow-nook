package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCatTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCatTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCatTasksLogic {
	return &ListCatTasksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 正式任务列表
func (l *ListCatTasksLogic) ListCatTasks(in *v1.ListCatTasksRequest) (*v1.ListCatTasksResponse, error) {
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	lastID := decodeCursor(in.Cursor)

	tasks, hasMore, err := l.svcCtx.Repo.Cat.ListCatTasks(l.ctx, lastID, pageSize)
	if err != nil {
		logger.Error("ListCatTasks: query failed")
		return nil, errorx.WrapInternal("查询任务列表失败", err)
	}

	items := make([]*v1.CatTaskItem, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, &v1.CatTaskItem{
			TaskId:          t.ID,
			CatId:           t.CatID,
			Title:           t.Title,
			TaskType:        t.TaskType,
			UrgencyLevel:    t.UrgencyLevel,
			DifficultyLevel: t.DifficultyLevel,
			RewardPoints:    t.RewardPoints,
			Status:          t.Status,
			Deadline:        timestamppb.New(t.Deadline),
			CreatedAt:       timestamppb.New(t.CreatedAt),
		})
	}

	nextCursor := ""
	if hasMore && len(tasks) > 0 {
		nextCursor = encodeCursor(tasks[len(tasks)-1].ID)
	}

	return &v1.ListCatTasksResponse{
		List:       items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
