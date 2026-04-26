package cat

import (
	"context"
	"time"

	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

// 从 JWT Context 取当前用户 ID（需和 Auth middleware 约定 key）
func getUserID(ctx context.Context) uint64 {
	uid, _ := ctx.Value("userId").(uint64)
	return uid
}

// 从 JWT Context 取用户角色
func getRole(ctx context.Context) string {
	role, _ := ctx.Value("role").(string)
	return role
}

func isAdmin(ctx context.Context) bool {
	return getRole(ctx) == "admin"
}

func convertImagesToPB(items []types.ImageItem) []*catpb.ImageItem {
	if len(items) == 0 {
		return nil
	}
	res := make([]*catpb.ImageItem, 0, len(items))
	for _, item := range items {
		res = append(res, &catpb.ImageItem{
			Url:         item.Url,
			Description: item.Description,
			Sort:        item.Sort,
			IsCover:     item.IsCover,
		})
	}
	return res
}

func convertImagesFromPB(items []*catpb.ImageItem) []types.ImageItem {
	if len(items) == 0 {
		return nil
	}
	res := make([]types.ImageItem, 0, len(items))
	for _, item := range items {
		res = append(res, types.ImageItem{
			Url:         item.Url,
			Description: item.Description,
			Sort:        item.Sort,
			IsCover:     item.IsCover,
		})
	}
	return res
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
