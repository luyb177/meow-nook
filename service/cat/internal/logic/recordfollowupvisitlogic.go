package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type RecordFollowUpVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecordFollowUpVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordFollowUpVisitLogic {
	return &RecordFollowUpVisitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecordFollowUpVisitLogic) RecordFollowUpVisit(in *v1.RecordFollowUpVisitRequest) (*v1.RecordFollowUpVisitResponse, error) {
	if in.AdoptionId == 0 || in.VisitorId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "adoption_id and visitor_id are required", errorx.ErrBadRequest)
	}
	if in.VisitType < 1 || in.VisitType > 4 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "invalid visit_type", errorx.ErrBadRequest)
	}

	if err := l.svcCtx.Repo.Adoption.RecordFollowUpVisit(
		l.ctx,
		in.AdoptionId,
		int(in.VisitType),
		in.VisitorId,
		time.Now(),
		in.Photos,
		in.Remark,
	); err != nil {
		return nil, errorx.WrapDBUpdate("记录回访失败", err)
	}

	return &v1.RecordFollowUpVisitResponse{
		Message: "记录回访成功",
	}, nil
}
