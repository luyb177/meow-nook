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

type ListCreateCatTaskApplyLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewListCreateCatTaskApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCreateCatTaskApplyLogic {
return &ListCreateCatTaskApplyLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

func (l *ListCreateCatTaskApplyLogic) ListCreateCatTaskApply(in *v1.ListCreateCatTaskApplyRequest) (*v1.ListCreateCatTaskApplyResponse, error) {
pageSize := int(in.PageSize)
if pageSize <= 0 {
pageSize = 20
}
page := int(in.Page)
if page <= 0 {
page = 1
}

applies, total, err := l.svcCtx.Repo.Task.ListTaskApplies(l.ctx, &task.TaskApplyFilter{
Status:   task.TaskApplyStatusPending,
Page:     page,
PageSize: pageSize,
})
if err != nil {
return nil, errorx.WrapDBQuery("查询任务申请列表失败", err)
}

items := make([]*v1.CreateCatTaskApplyItem, 0, len(applies))
for _, a := range applies {
items = append(items, &v1.CreateCatTaskApplyItem{
ApplyId:         a.ID,
CatId:           a.CatID,
Title:           a.Title,
TaskType:        a.TaskType,
ApplicantUserId: a.ApplicantUserID,
Status:          a.Status,
CreatedAt:       timestamppb.New(a.CreatedAt),
})
}

return &v1.ListCreateCatTaskApplyResponse{
List:     items,
Total:    total,
Page:     int32(page),
PageSize: int32(pageSize),
}, nil
}
