package verify

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	v1 "github.com/luyb177/meow-nook/service/user/pb/user/v1"
	"github.com/redis/go-redis/v9"
)

const (
	// verify:code:{channel}:{purpose}:{target_hash}
	CodeKey = "verify:code:%d:%d:%s"
)

type VerifyMeta struct {
	Target  string
	Channel v1.VerifyChannel
	Purpose v1.VerifyPurpose
}

type Repository interface {
	// SetCode 设置验证码
	SetCode(ctx context.Context, meta *VerifyMeta, code string, expire time.Duration) error

	// GetAndDeleteCode 原子获取并删除验证码（Lua 保证一致性）
	GetAndDeleteCode(ctx context.Context, meta *VerifyMeta) (string, bool, error)

	// SetVerified 在验证码校验通过后打一个"已验证"标记，供后续业务（如注册）消费。
	// expire 建议 10 分钟，业务消费后应调用 DeleteVerified 清除。
	SetVerified(ctx context.Context, meta *VerifyMeta, expire time.Duration) error

	// GetAndDeleteVerified 原子读取并删除"已验证"标记，返回是否存在。
	GetAndDeleteVerified(ctx context.Context, meta *VerifyMeta) (bool, error)
}

type repo struct {
	client *redis.Client
}

func NewVerifyRepo(client *redis.Client) Repository {
	return &repo{
		client: client,
	}
}

func (r *repo) SetCode(ctx context.Context, meta *VerifyMeta, code string, expire time.Duration) error {
	key := verifyCodeKey(meta)

	return r.client.Set(ctx, key, code, expire).Err()
}

func (r *repo) GetAndDeleteCode(ctx context.Context, meta *VerifyMeta) (string, bool, error) {
	key := verifyCodeKey(meta)

	val, err := getAndDeleteScript.Run(
		ctx,
		r.client,
		[]string{key},
	).Result()

	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}

	if err != nil {
		return "", false, err
	}

	str, ok := val.(string)
	if !ok {
		return "", false, fmt.Errorf("unexpected redis result type: %T", val)
	}

	return str, true, nil
}

func verifyCodeKey(meta *VerifyMeta) string {
	// target 做 hash，避免邮箱/手机号直接暴露在 key 中
	sum := sha256.Sum256([]byte(meta.Target))
	targetHash := hex.EncodeToString(sum[:])

	return fmt.Sprintf(CodeKey, meta.Channel, meta.Purpose, targetHash)
}

// ──────────────────────────────────────────────
// Verified mark
// ──────────────────────────────────────────────

const (
	// verify:verified:{channel}:{purpose}:{target_hash}
	verifiedKey = "verify:verified:%d:%d:%s"
)

func (r *repo) SetVerified(ctx context.Context, meta *VerifyMeta, expire time.Duration) error {
	key := verifiedMarkKey(meta)
	return r.client.Set(ctx, key, "1", expire).Err()
}

func (r *repo) GetAndDeleteVerified(ctx context.Context, meta *VerifyMeta) (bool, error) {
	key := verifiedMarkKey(meta)

	val, err := getAndDeleteScript.Run(ctx, r.client, []string{key}).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val != nil, nil
}

func verifiedMarkKey(meta *VerifyMeta) string {
	sum := sha256.Sum256([]byte(meta.Target))
	targetHash := hex.EncodeToString(sum[:])
	return fmt.Sprintf(verifiedKey, meta.Channel, meta.Purpose, targetHash)
}
