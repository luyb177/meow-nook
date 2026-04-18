package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"
)

type SendVerificationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendVerificationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerificationCodeLogic {
	return &SendVerificationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SendVerificationCode Auth
func (l *SendVerificationCodeLogic) SendVerificationCode(in *v1.SendVerificationCodeReq) (*v1.Response, error) {
	logger.Info("User SendVerificationCode")

	switch in.GetChannel() {
	case v1.VerifyChannel_VERIFY_CHANNEL_EMAIL:
		return l.sendVerificationCodeByEmail(in)
	case v1.VerifyChannel_VERIFY_CHANNEL_PHONE:
		return nil, errorx.ErrNotImplemented
	default:
		return nil, errorx.ErrBadRequest
	}
}

func (l *SendVerificationCodeLogic) sendVerificationCodeByEmail(in *v1.SendVerificationCodeReq) (*v1.Response, error) {
	return &v1.Response{}, nil
}
