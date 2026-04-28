package logic

import (
	"context"

	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type CompleteVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteVisitLogic {
	return &CompleteVisitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CompleteVisitLogic) CompleteVisit(in *v1.RecordFollowUpVisitRequest) (*v1.RecordFollowUpVisitResponse, error) {
	recordLogic := NewRecordFollowUpVisitLogic(l.ctx, l.svcCtx)
	return recordLogic.RecordFollowUpVisit(in)
}
