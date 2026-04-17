package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFeedbackLogic {
	return &SubmitFeedbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Feedback
func (l *SubmitFeedbackLogic) SubmitFeedback(in *v1.SubmitFeedbackReq) (*v1.SubmitFeedbackResp, error) {
	// todo: add your logic here and delete this line

	return &v1.SubmitFeedbackResp{}, nil
}
