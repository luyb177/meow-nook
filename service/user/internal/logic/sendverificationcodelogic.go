package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/code"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/email"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/task"
	"github.com/luyb177/meow-nook/service/user/internal/repo/verify"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendVerificationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendVerificationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerificationCodeLogic {
	return &SendVerificationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendVerificationCodeLogic) SendVerificationCode(in *v1.SendVerificationCodeReq) (*v1.Response, error) {
	logger.Info("user SendVerificationCode called")

	if !validVerifyPurpose(in.Purpose) {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "不支持的验证码发送目的", errorx.ErrBadRequest)
	}

	switch in.Channel {
	case v1.VerifyChannel_VERIFY_CHANNEL_EMAIL:
		return l.sendVerificationEmailCode(in.Target, in.Purpose)
	default:
		return nil, errorx.Wrap(errorx.CodeBadRequest, "不支持的验证码发送渠道", errorx.ErrBadRequest)
	}
}

func (l *SendVerificationCodeLogic) sendVerificationEmailCode(target string, purpose v1.VerifyPurpose) (*v1.Response, error) {
	if target == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "邮箱地址不能为空", errorx.ErrBadRequest)
	}

	if len(target) > 254 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "邮箱地址长度不能超过254个字符", errorx.ErrBadRequest)
	}

	target = email.CanonicalEmail(target)

	emailCode := code.EmailCode()

	err := l.svcCtx.Repo.Verify.SetCode(l.ctx, &verify.VerifyMeta{
		Target:  target,
		Channel: v1.VerifyChannel_VERIFY_CHANNEL_EMAIL,
		Purpose: purpose,
	}, emailCode, time.Minute*5)

	if err != nil {
		return nil, errorx.WrapInternal("设置验证码失败", err)
	}

	ctx, cancel := context.WithTimeout(l.ctx, 5*time.Second)
	defer cancel()
	t1 := kafka.NewTypedTask(task.TypeSendVerifyCode, "send-verify-code", task.VerifyCode{To: target, Code: emailCode})
	if err = l.svcCtx.KafkaProducer.Dispatch(ctx, t1); err != nil {
		return nil, errorx.WrapInternal("投递<发送验证码>任务失败", err)
	}
	return &v1.Response{}, nil
}

func validVerifyPurpose(purpose v1.VerifyPurpose) bool {
	switch purpose {
	case v1.VerifyPurpose_VERIFY_PURPOSE_LOGIN,
		v1.VerifyPurpose_VERIFY_PURPOSE_REGISTER,
		v1.VerifyPurpose_VERIFY_PURPOSE_RESET_PASSWORD:
		return true
	default:
		return false
	}
}
