package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	pkgemail "github.com/luyb177/meow-nook/service/user/internal/pkg/email"
	"github.com/luyb177/meow-nook/service/user/internal/repo/verify"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCodeLogic {
	return &VerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyCodeLogic) VerifyCode(in *v1.VerifyCodeReq) (*v1.Response, error) {
	logger.Info("user VerifyCode called")

	if in.Target == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "目标地址不能为空", errorx.ErrBadRequest)
	}
	if in.Code == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "验证码不能为空", errorx.ErrBadRequest)
	}
	if !validVerifyPurpose(in.Purpose) {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "不支持的验证码用途", errorx.ErrBadRequest)
	}

	target := pkgemail.CanonicalEmail(in.Target)

	meta := &verify.VerifyMeta{
		Target:  target,
		Channel: in.Channel,
		Purpose: in.Purpose,
	}

	// 原子读取并删除验证码
	stored, exists, err := l.svcCtx.Repo.Verify.GetAndDeleteCode(l.ctx, meta)
	if err != nil {
		return nil, errorx.WrapInternal("读取验证码失败", err)
	}
	if !exists {
		return nil, errorx.New(errorx.CodeBadRequest, "验证码已过期或不存在")
	}
	if stored != in.Code {
		return nil, errorx.New(errorx.CodeBadRequest, "验证码错误")
	}

	// 打"已验证"标记，有效期 10 分钟，供后续注册/重置密码消费
	if err = l.svcCtx.Repo.Verify.SetVerified(l.ctx, meta, 10*time.Minute); err != nil {
		return nil, errorx.WrapInternal("设置验证标记失败", err)
	}

	return &v1.Response{}, nil
}
