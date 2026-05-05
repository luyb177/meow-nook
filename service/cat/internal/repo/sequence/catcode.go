package sequence

import (
	"context"
	"fmt"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
)

// GenerateCatCode 生成猫咪编号
//
// 格式：
// CAT-20260426-000001
func (r *repository) GenerateCatCode(ctx context.Context) (string, error) {
	date := time.Now().Format("20060102")
	key := fmt.Sprintf("cat_code:%s", date)

	// Redis 原子自增
	seq, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return "", errorx.WrapRedisSet("生成猫咪编号失败", err)
	}

	// 设置过期时间（可选）
	// 防止 Redis 永久堆积历史 key
	_ = r.client.Expire(ctx, key, 48*time.Hour).Err()

	return fmt.Sprintf("CAT-%s-%06d", date, seq), nil
}
