package logic

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v4"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	pkgemail "github.com/luyb177/meow-nook/service/user/internal/pkg/email"
	usermodel "github.com/luyb177/meow-nook/service/user/internal/repo/user"
	"github.com/luyb177/meow-nook/service/user/internal/repo/verify"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *v1.RegisterReq) (*v1.RegisterResp, error) {
	logger.Info("user Register called")

	// ── 1. 参数校验 ──────────────────────────────────────────────
	if in.Email == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "邮箱不能为空", errorx.ErrBadRequest)
	}
	if !pkgemail.IsValidEmail(in.Email) {
		return nil, errorx.New(errorx.CodeBadRequest, "邮箱格式不正确")
	}
	if in.Password == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "密码不能为空", errorx.ErrBadRequest)
	}
	if utf8.RuneCountInString(in.Password) < 8 {
		return nil, errorx.New(errorx.CodeBadRequest, "密码长度不能少于8位")
	}

	email := pkgemail.CanonicalEmail(in.Email)

	// ── 2. 检查"已验证"标记（由 VerifyCode 接口写入）──────────────
	meta := &verify.VerifyMeta{
		Target:  email,
		Channel: v1.VerifyChannel_VERIFY_CHANNEL_EMAIL,
		Purpose: v1.VerifyPurpose_VERIFY_PURPOSE_REGISTER,
	}
	verified, err := l.svcCtx.Repo.Verify.GetAndDeleteVerified(l.ctx, meta)
	if err != nil {
		return nil, errorx.WrapInternal("读取验证标记失败", err)
	}
	if !verified {
		return nil, errorx.New(errorx.CodeBadRequest, "邮箱尚未完成验证，请先验证验证码")
	}

	// ── 3. 检查邮箱是否已注册 ────────────────────────────────────
	exists, err := l.svcCtx.Repo.User.ExistsByEmail(l.ctx, email)
	if err != nil {
		return nil, errorx.WrapInternal("查询用户失败", err)
	}
	if exists {
		return nil, errorx.ErrUserAlreadyExists
	}

	// ── 4. 密码哈希（salt + sha256，无需外部依赖）────────────────
	passwordHash, err := hashPassword(in.Password)
	if err != nil {
		return nil, errorx.WrapInternal("密码加密失败", err)
	}

	// ── 5. 入库 ──────────────────────────────────────────────────
	u := &usermodel.User{
		Email:        email,
		PasswordHash: passwordHash,
		Username:     email, // 默认用邮箱作为用户名，后续可修改
		Role:         "user",
	}
	if err = l.svcCtx.Repo.User.Create(l.ctx, u); err != nil {
		return nil, errorx.WrapInternal("创建用户失败", err)
	}

	// ── 6. 注册即登录，签发 JWT（包含 user_id + role，有效期 14 天）
	token, err := generateJWT(l.svcCtx.Config.JWT.Secret, l.svcCtx.Config.JWT.ExpireTime, u.ID, u.Role)
	if err != nil {
		return nil, errorx.WrapInternal("生成 token 失败", err)
	}

	return &v1.RegisterResp{UserId: u.ID, Token: token}, nil
}

// hashPassword 生成 "salt$hash" 格式的密码摘要。
func hashPassword(password string) (string, error) {
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", err
	}
	salt := hex.EncodeToString(saltBytes)
	sum := sha256.Sum256([]byte(salt + password))
	return salt + "$" + hex.EncodeToString(sum[:]), nil
}

// generateJWT 是跨 logic 共用的 JWT 签发函数。
func generateJWT(secret string, expireTime time.Duration, userID int64, role string) (string, error) {
	expire := expireTime
	if expire <= 0 {
		expire = 14 * 24 * time.Hour
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(expire).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
