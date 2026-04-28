package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type RecordHomeVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecordHomeVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordHomeVisitLogic {
	return &RecordHomeVisitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecordHomeVisitLogic) RecordHomeVisit(in *v1.RecordHomeVisitRequest) (*v1.RecordHomeVisitResponse, error) {
	if in.AdoptionId == 0 || in.VisitorId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "adoption_id and visitor_id are required", errorx.ErrBadRequest)
	}

	if err := l.svcCtx.Repo.Adoption.RecordHomeVisit(
		l.ctx,
		in.AdoptionId,
		time.Now(),
		in.VisitorId,
		in.Photos,
		in.Remark,
	); err != nil {
		return nil, errorx.WrapDBUpdate("记录家访失败", err)
	}

	return &v1.RecordHomeVisitResponse{
		Message: "记录家访成功",
	}, nil
}
