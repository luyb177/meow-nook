package logic

import (
"context"

"github.com/luyb177/meow-nook/common/errorx"
"github.com/luyb177/meow-nook/service/cat/internal/repo/task"
"github.com/luyb177/meow-nook/service/cat/internal/svc"
v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
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
page := int(in.Page)
if page <= 0 {
page = 1
}

filter := &task.TaskFilter{
Status:       in.Status,
UrgencyLevel: in.UrgencyLevel,
Area:         in.Area,
Page:         page,
PageSize:     pageSize,
}
if in.CatId > 0 {
filter.CatID = in.CatId
}
if in.DifficultyLevel > 0 {
filter.DifficultyLevel = in.DifficultyLevel
}

tasks, total, err := l.svcCtx.Repo.Task.ListTasks(l.ctx, filter)
if err != nil {
return nil, errorx.WrapDBQuery("查询任务列表失败", err)
}

items := make([]*v1.CatTaskItem, 0, len(tasks))
for _, t := range tasks {
item := &v1.CatTaskItem{
TaskId:          t.ID,
Title:           t.Title,
TaskType:        t.TaskType,
UrgencyLevel:    t.UrgencyLevel,
DifficultyLevel: t.DifficultyLevel,
RewardPoints:    t.RewardPoints,
Status:          t.Status,
CreatedAt:       timestamppb.New(t.CreatedAt),
}
if t.DeadlineAt != nil {
item.Deadline = timestamppb.New(*t.DeadlineAt)
}
items = append(items, item)
}

return &v1.ListCatTasksResponse{
List:     items,
Total:    total,
Page:     int32(page),
PageSize: int32(pageSize),
}, nil
}
