package logic

import (
"context"

"github.com/luyb177/meow-nook/common/errorx"
"github.com/luyb177/meow-nook/service/cat/internal/svc"
v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
"google.golang.org/protobuf/types/known/timestamppb"

"github.com/zeromicro/go-zero/core/logx"
)

type GetCatTaskDetailLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewGetCatTaskDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCatTaskDetailLogic {
return &GetCatTaskDetailLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 获取任务详情
func (l *GetCatTaskDetailLogic) GetCatTaskDetail(in *v1.GetCatTaskDetailRequest) (*v1.GetCatTaskDetailResponse, error) {
if in.TaskId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "task_id不能为空", errorx.ErrBadRequest)
}

t, err := l.svcCtx.Repo.Task.GetTaskByID(l.ctx, in.TaskId)
if err != nil {
return nil, errorx.WrapDBQuery("查询任务失败", err)
}

resp := &v1.GetCatTaskDetailResponse{
TaskId:          t.ID,
CatId:           t.CatID,
Title:           t.Title,
TaskType:        t.TaskType,
UrgencyLevel:    t.UrgencyLevel,
DifficultyLevel: t.DifficultyLevel,
RewardPoints:    t.RewardPoints,
Summary:         t.Summary,
Detail:          t.Detail,
Status:          t.Status,
Remark:          t.Remark,
MaxClaimers:     t.MaxClaimers,
CurrentClaimers: t.CurrentClaimers,
CreatedAt:       timestamppb.New(t.CreatedAt),
UpdatedAt:       timestamppb.New(t.UpdatedAt),
}
if t.DeadlineAt != nil {
resp.Deadline = timestamppb.New(*t.DeadlineAt)
}

return resp, nil
}
