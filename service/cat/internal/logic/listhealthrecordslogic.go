package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListHealthRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListHealthRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListHealthRecordsLogic {
	return &ListHealthRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListHealthRecordsLogic) ListHealthRecords(in *v1.ListHealthRecordsReq) (*v1.ListHealthRecordsResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ListHealthRecordsResp{}, nil
}
