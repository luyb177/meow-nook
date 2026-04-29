package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetTaskFlowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskFlowsLogic {
	return &GetTaskFlowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskFlowsLogic) GetTaskFlows(in *v1.GetTaskFlowsRequest) (*v1.GetTaskFlowsResponse, error) {
	logger.Info("GetTaskFlows called")

	flows, err := l.svcCtx.Repo.Task.ListTaskFlows(l.ctx, in.TaskId)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询任务流程失败", err)
	}

	items := make([]*v1.TaskFlowVO, 0, len(flows))
	for _, v := range flows {
		items = append(items, &v1.TaskFlowVO{
			Id:         v.ID,
			TaskId:     v.TaskID,
			UserId:     v.UserID,
			Action:     v.Action,
			FromStatus: v.FromStatus,
			ToStatus:   v.ToStatus,
			Remark:     v.Remark,
			CreatedAt:  timestamppb.New(v.CreatedAt),
		})
	}

	return &v1.GetTaskFlowsResponse{Items: items}, nil
}
