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

type GetMyTaskApplyDetailLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewGetMyTaskApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyTaskApplyDetailLogic {
return &GetMyTaskApplyDetailLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// 志愿者获取自己的任务申请详情
func (l *GetMyTaskApplyDetailLogic) GetMyTaskApplyDetail(in *v1.GetMyTaskApplyDetailRequest) (*v1.GetMyTaskApplyDetailResponse, error) {
if in.ApplyId == 0 {
return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id不能为空", errorx.ErrBadRequest)
}

apply, err := l.svcCtx.Repo.Task.GetTaskApplyByID(l.ctx, in.ApplyId)
if err != nil {
return nil, errorx.WrapDBQuery("查询任务申请失败", err)
}

// 权限校验：只能查自己的申请
if in.UserId > 0 && apply.ApplicantUserID != in.UserId {
return nil, errorx.Wrap(errorx.CodeForbidden, "无权查看此申请", errorx.ErrForbidden)
}

resp := &v1.GetMyTaskApplyDetailResponse{
ApplyId:         apply.ID,
CatId:           apply.CatID,
Title:           apply.Title,
TaskType:        apply.TaskType,
Summary:         apply.Summary,
Detail:          apply.Detail,
Status:          apply.Status,
RejectReason:    apply.RejectReason,
UrgencyLevel:    apply.UrgencyLevel,
DifficultyLevel: apply.DifficultyLevel,
RewardPoints:    apply.RewardPoints,
CreatedAt:       timestamppb.New(apply.CreatedAt),
UpdatedAt:       timestamppb.New(apply.UpdatedAt),
}
if apply.DeadlineAt != nil {
resp.Deadline = timestamppb.New(*apply.DeadlineAt)
}

_ = task.TaskApplyStatusPending // ensure package used

return resp, nil
}
