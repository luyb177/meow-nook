package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyAddRescueRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyAddRescueRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAddRescueRecordLogic {
	return &ApplyAddRescueRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请新增救助记录
func (l *ApplyAddRescueRecordLogic) ApplyAddRescueRecord(in *v1.ApplyAddRescueRecordRequest) (*v1.ApplyAddRescueRecordResponse, error) {
	if _, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, in.CatId); err != nil {
		return nil, errorx.ErrCatNotFound
	}

	var rescueTime = in.RescueTime.AsTime()
	apply := &cat.CatRescueApply{
		CatID:           in.CatId,
		RescueTime:      rescueTime,
		RescuerName:     in.RescuerName,
		RescueStatus:    in.RescueStatus,
		Description:     in.Description,
		ApplicantUserID: in.OperatorId,
		Status:          "pending",
	}

	if err := l.svcCtx.Repo.Cat.CreateCatRescueApply(l.ctx, apply); err != nil {
		logger.Error("ApplyAddRescueRecord: create apply failed")
		return nil, errorx.WrapInternal("创建救助记录申请失败", err)
	}

	return &v1.ApplyAddRescueRecordResponse{RescueApplyId: apply.ID}, nil
}
