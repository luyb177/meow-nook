package tag

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	// ==================== Tag 定义相关 ====================

	// CreateTag 创建标签
	CreateTag(ctx context.Context, tag *Tag, tx ...*gorm.DB) error

	// GetTagByID 获取标签详情
	GetTagByID(ctx context.Context, tagID uint64, tx ...*gorm.DB) (*Tag, error)

	// GetTagByNameAndType 按名称和类型查询（用于防止重复创建）
	GetTagByNameAndType(ctx context.Context, name, tagType string, tx ...*gorm.DB) (*Tag, error)

	// ListTags 标签列表
	ListTags(ctx context.Context, filter TagListFilter, tx ...*gorm.DB) ([]*Tag, int64, error)

	// UpdateTag 更新标签（不允许修改 name + type 组合）
	UpdateTag(ctx context.Context, tagID uint64, values map[string]any, tx ...*gorm.DB) error

	// DeleteTag 软删除标签（同时删除关联关系）
	DeleteTag(ctx context.Context, tagID uint64, tx ...*gorm.DB) error

	// ==================== Tag 关联相关 ====================

	// AddTagRelation 添加标签关联
	AddTagRelation(ctx context.Context, tagID uint64, targetType string, targetID uint64, sort int32, tx ...*gorm.DB) error

	// BatchAddTagRelations 批量添加标签关联
	BatchAddTagRelations(ctx context.Context, relations []*TagRelation, tx ...*gorm.DB) error

	// RemoveTagRelation 删除标签关联
	RemoveTagRelation(ctx context.Context, tagID uint64, targetType string, targetID uint64, tx ...*gorm.DB) error

	// ListTagRelationsByTarget 查询某个对象的所有标签
	ListTagRelationsByTarget(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) ([]*TagRelation, error)

	// GetTagsByTarget 获取某个对象的标签详情（带 Tag 信息）
	GetTagsByTarget(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) ([]*Tag, error)

	// ReplaceTagsByTarget 替换某个对象的所有标签（先删除旧的，再添加新的）
	ReplaceTagsByTarget(ctx context.Context, targetType string, targetID uint64, tagIDs []uint64, tx ...*gorm.DB) error
}

type repository struct {
	db     *gorm.DB
	client *redis.Client
}

func NewRepository(db *gorm.DB, client *redis.Client) Repository {
	return &repository{
		db:     db,
		client: client,
	}
}
