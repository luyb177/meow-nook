package logic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	pkgemail "github.com/luyb177/meow-nook/service/user/internal/pkg/email"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *v1.LoginReq) (*v1.LoginResp, error) {
	logger.Info("user Login called")

	// ── 1. 参数校验 ──────────────────────────────────────────────
	if in.Email == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "邮箱不能为空", errorx.ErrBadRequest)
	}
	if in.Password == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "密码不能为空", errorx.ErrBadRequest)
	}

	email := pkgemail.CanonicalEmail(in.Email)

	// ── 2. 查询用户是否存在 ──────────────────────────────────────
	u, err := l.svcCtx.Repo.User.FindByEmail(l.ctx, email)
	if err != nil {
		return nil, errorx.WrapInternal("查询用户失败", err)
	}
	if u == nil {
		return nil, errorx.ErrUserNotFound
	}

	// ── 3. 校验密码 ──────────────────────────────────────────────
	if !verifyPassword(in.Password, u.PasswordHash) {
		return nil, errorx.ErrPasswordWrong
	}

	// ── 4. 签发 JWT（包含 user_id + role，有效期 14 天）──────────
	token, err := generateJWT(l.svcCtx.Config.JWT.Secret, l.svcCtx.Config.JWT.ExpireTime, u.ID, u.Role)
	if err != nil {
		return nil, errorx.WrapInternal("生成 token 失败", err)
	}

	return &v1.LoginResp{
		Token:  token,
		UserId: u.ID,
	}, nil
}

// verifyPassword 验证明文密码与存储的 "salt$hash" 是否匹配。
func verifyPassword(password, stored string) bool {
	parts := strings.SplitN(stored, "$", 2)
	if len(parts) != 2 {
		return false
	}
	salt := parts[0]
	sum := sha256.Sum256([]byte(salt + password))
	return hex.EncodeToString(sum[:]) == parts[1]
}

// generateToken 签发 JWT，claims 包含 user_id、role，有效期 14 天。
func (l *LoginLogic) generateToken(userID int64, role string) (string, error) {
	cfg := l.svcCtx.Config.JWT
	expire := cfg.ExpireTime
	if expire <= 0 {
		expire = 14 * 24 * time.Hour // 默认 14 天
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(expire).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}
