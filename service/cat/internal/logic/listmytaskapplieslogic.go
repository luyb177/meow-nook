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

type ListMyTaskAppliesLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewListMyTaskAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyTaskAppliesLogic {
return &ListMyTaskAppliesLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 我的任务申请列表
func (l *ListMyTaskAppliesLogic) ListMyTaskApplies(in *v1.ListMyTaskAppliesRequest) (*v1.ListMyTaskAppliesResponse, error) {
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

applies, total, err := l.svcCtx.Repo.Task.ListTaskApplies(l.ctx, &task.TaskApplyFilter{
ApplicantUserID: in.UserId,
Status:          in.Status,
Page:            page,
PageSize:        pageSize,
})
if err != nil {
return nil, errorx.WrapDBQuery("查询我的任务申请失败", err)
}

items := make([]*v1.TaskApplyItem, 0, len(applies))
for _, a := range applies {
items = append(items, &v1.TaskApplyItem{
ApplyId:   a.ID,
CatId:     a.CatID,
Title:     a.Title,
TaskType:  a.TaskType,
Status:    a.Status,
CreatedAt: timestamppb.New(a.CreatedAt),
})
}

return &v1.ListMyTaskAppliesResponse{
List:     items,
Total:    total,
Page:     int32(page),
PageSize: int32(pageSize),
}, nil
}
