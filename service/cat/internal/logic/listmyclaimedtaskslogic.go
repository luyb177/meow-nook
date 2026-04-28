package logic

import (
"context"

"github.com/luyb177/meow-nook/common/errorx"
"github.com/luyb177/meow-nook/service/cat/internal/svc"
v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
"google.golang.org/protobuf/types/known/timestamppb"

"github.com/zeromicro/go-zero/core/logx"
)

type ListMyClaimedTasksLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewListMyClaimedTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyClaimedTasksLogic {
return &ListMyClaimedTasksLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 获取我认领的任务列表
func (l *ListMyClaimedTasksLogic) ListMyClaimedTasks(in *v1.ListMyClaimedTasksRequest) (*v1.ListMyClaimedTasksResponse, error) {
if in.UserId == 0 {
return nil, errorx.Wrap(errorx.CodeUnauthorized, "用户未登录", errorx.ErrUnauthorized)
}

page := int(in.Page)
if page <= 0 {
page = 1
}
pageSize := int(in.PageSize)
if pageSize <= 0 {
pageSize = 20
}

tasks, claims, total, err := l.svcCtx.Repo.Task.ListMyClaimedTasks(l.ctx, in.UserId, in.ClaimStatus, page, pageSize)
if err != nil {
return nil, errorx.WrapDBQuery("查询我认领的任务失败", err)
}

// Build a map from task_id to claim
claimMap := make(map[uint64]*struct {
status    string
claimedAt timestamppb.Timestamp
})
for _, c := range claims {
claimMap[c.TaskID] = &struct {
status    string
claimedAt timestamppb.Timestamp
}{
status:    c.Status,
claimedAt: *timestamppb.New(c.CreatedAt),
}
}

items := make([]*v1.ClaimedTaskItem, 0, len(tasks))
for _, t := range tasks {
item := &v1.ClaimedTaskItem{
TaskId:          t.ID,
CatId:           t.CatID,
Title:           t.Title,
TaskType:        t.TaskType,
UrgencyLevel:    t.UrgencyLevel,
DifficultyLevel: t.DifficultyLevel,
RewardPoints:    t.RewardPoints,
TaskStatus:      t.Status,
}
if cm, ok := claimMap[t.ID]; ok {
item.ClaimStatus = cm.status
item.ClaimedAt = &cm.claimedAt
}
if t.DeadlineAt != nil {
item.Deadline = timestamppb.New(*t.DeadlineAt)
}
items = append(items, item)
}

return &v1.ListMyClaimedTasksResponse{
List:     items,
Total:    total,
Page:     int32(page),
PageSize: int32(pageSize),
}, nil
}
