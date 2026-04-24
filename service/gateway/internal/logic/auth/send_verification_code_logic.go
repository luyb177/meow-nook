// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	userpb "github.com/luyb177/meow-nook/service/user/pb/user/v1"
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

func (l *SendVerificationCodeLogic) SendVerificationCode(req *types.SendVerificationCodeReq) (resp *types.Response, err error) {
	logger.Info("SendVerificationCodeLogic called")

	if !validVerifyChannel(req.Channel) {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "不支持的验证码发送渠道", errorx.ErrBadRequest)
	}

	if !validVerifyPurpose(req.Purpose) {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "不支持的验证码发送目的", errorx.ErrBadRequest)
	}

	_, err = l.svcCtx.UserRPC.SendVerificationCode(l.ctx, &userpb.SendVerificationCodeReq{
		Target:  req.Target,
		Channel: userpb.VerifyChannel(req.Channel),
		Purpose: userpb.VerifyPurpose(req.Purpose),
	})

	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.Response{}, nil
}

func validVerifyChannel(channel int32) bool {
	switch userpb.VerifyChannel(channel) {
	case userpb.VerifyChannel_VERIFY_CHANNEL_EMAIL,
		userpb.VerifyChannel_VERIFY_CHANNEL_PHONE:
		return true
	default:
		return false
	}
}

func validVerifyPurpose(purpose int32) bool {
	switch userpb.VerifyPurpose(purpose) {
	case userpb.VerifyPurpose_VERIFY_PURPOSE_LOGIN,
		userpb.VerifyPurpose_VERIFY_PURPOSE_REGISTER,
		userpb.VerifyPurpose_VERIFY_PURPOSE_RESET_PASSWORD:
		return true
	default:
		return false
	}
}
