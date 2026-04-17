package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRescueRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRescueRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRescueRecordsLogic {
	return &ListRescueRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListRescueRecordsLogic) ListRescueRecords(in *v1.ListRescueRecordsReq) (*v1.ListRescueRecordsResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ListRescueRecordsResp{}, nil
}
