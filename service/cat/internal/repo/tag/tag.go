package tag

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// ==================== Tag 定义表 ====================

type Tag struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	Name        string `gorm:"type:varchar(64);uniqueIndex;not null;comment:标签名称" json:"name"`
	Type        string `gorm:"type:varchar(32);not null;index;comment:标签类型 cat/post/task/user" json:"type"`
	Theme       string `gorm:"type:varchar(32);comment:主题色 success/warning/danger/info" json:"theme"`
	Priority    int32  `gorm:"default:0;comment:优先级，越大越靠前" json:"priority"`
	Description string `gorm:"type:varchar(255);comment:标签描述" json:"description"`
	CreatedBy   uint64 `gorm:"default:0;comment:创建人ID" json:"created_by"`

	CreatedAt time.Time             `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano;comment:删除时间" json:"deleted_at,omitempty"`
}

func (Tag) TableName() string {
	return "tags"
}

// ==================== Tag 关联表 ====================

type TagRelation struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	TagID      uint64 `gorm:"index;not null;comment:标签ID" json:"tag_id"`
	TargetID   uint64 `gorm:"index;not null;comment:关联对象ID" json:"target_id"`
	TargetType string `gorm:"type:varchar(32);index;not null;comment:关联对象类型 cat/post/task/user" json:"target_type"`
	Sort       int32  `gorm:"default:0;comment:排序，越大越靠前" json:"sort"`
	CreatedBy  uint64 `gorm:"default:0;comment:创建人ID" json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TagRelation) TableName() string {
	return "tag_relations"
}

const (
	// 标签类型
	TagTypeCat  = "cat"  // 猫咪标签
	TagTypePost = "post" // 动态标签
	TagTypeTask = "task" // 任务标签
	TagTypeUser = "user" // 用户标签

	// 主题色
	TagThemeSuccess = "success" // 绿色
	TagThemeWarning = "warning" // 黄色
	TagThemeDanger  = "danger"  // 红色
	TagThemeInfo    = "info"    // 蓝色
)

// 猫咪标签类型示例
const (
	TagCatTypePersonality = "personality" // 性格：活泼、亲人、高冷
	TagCatTypeHealth      = "health"      // 健康：已绝育、已疫苗、需医疗
	TagCatTypeAppearance  = "appearance"  // 外观：橘猫、长毛、白色
	TagCatTypeSource      = "source"      // 来源：救助、领养、流浪
)

const (
	TargetTypeCatApply   = "cat_apply"
	TargetTypeCatProfile = "cat_profile"

	// Tag Target Types
	TagTargetTypeCatTask      = "cat_task"
	TagTargetTypeCatTaskApply = "cat_task_apply"
	TagTargetTypeCatTaskClaim = "cat_task_claim"
	TagTargetTypeCatTaskLog   = "cat_task_log"
)

var (
	ErrTagNotFound       = errors.New("tag not found")
	ErrTagAlreadyExists  = errors.New("tag already exists")
	ErrTagRelationExists = errors.New("tag relation already exists")
	ErrInvalidTagType    = errors.New("invalid tag type")
	ErrInvalidTagTheme   = errors.New("invalid tag theme")
	ErrTagNameRequired   = errors.New("tag name is required")
)

func (r *repository) CreateTag(ctx context.Context, tag *Tag, tx ...*gorm.DB) error {
	if tag == nil {
		return errors.New("tag is nil")
	}
	if tag.Name == "" {
		return ErrTagNameRequired
	}
	if tag.Type == "" {
		return ErrInvalidTagType
	}

	// 检查是否已存在同名同类型的标签
	if existing, err := r.GetTagByNameAndType(ctx, tag.Name, tag.Type, tx...); err == nil && existing != nil {
		return ErrTagAlreadyExists
	}

	return r.getDB(ctx, tx...).Create(tag).Error
}

func (r *repository) AddTagRelation(ctx context.Context, tagID uint64, targetType string, targetID uint64, sort int32, tx ...*gorm.DB) error {
	if tagID == 0 {
		return ErrTagNotFound
	}
	if targetType == "" || targetID == 0 {
		return errors.New("invalid target")
	}

	// 检查是否已存在
	var count int64
	if err := r.getDB(ctx, tx...).
		Model(&TagRelation{}).
		Where("tag_id = ? AND target_type = ? AND target_id = ?", tagID, targetType, targetID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrTagRelationExists
	}

	relation := &TagRelation{
		TagID:      tagID,
		TargetID:   targetID,
		TargetType: targetType,
		Sort:       sort,
	}

	return r.getDB(ctx, tx...).Create(relation).Error
}

func (r *repository) BatchAddTagRelations(ctx context.Context, relations []*TagRelation, tx ...*gorm.DB) error {
	if len(relations) == 0 {
		return nil
	}

	return r.getDB(ctx, tx...).Create(relations).Error
}

func (r *repository) GetTagsByTarget(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) ([]*Tag, error) {
	var tags []*Tag
	err := r.getDB(ctx, tx...).
		Table("tags").
		Joins("JOIN tag_relations ON tags.id = tag_relations.tag_id").
		Where("tag_relations.target_type = ? AND tag_relations.target_id = ?", targetType, targetID).
		Where("tags.deleted_at IS NULL").
		Order("tag_relations.sort DESC, tags.priority DESC").
		Find(&tags).Error

	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *repository) ReplaceTagsByTarget(ctx context.Context, targetType string, targetID uint64, tagIDs []uint64, tx ...*gorm.DB) error {
	if len(tagIDs) == 0 {
		// 删除所有旧标签
		return r.getDB(ctx, tx...).
			Where("target_type = ? AND target_id = ?", targetType, targetID).
			Delete(&TagRelation{}).Error
	}

	return r.getDB(ctx, tx...).Transaction(func(tx *gorm.DB) error {
		// 1. 删除旧标签
		if err := tx.
			Where("target_type = ? AND target_id = ?", targetType, targetID).
			Delete(&TagRelation{}).Error; err != nil {
			return err
		}

		// 2. 添加新标签
		relations := make([]*TagRelation, 0, len(tagIDs))
		for i, tagID := range tagIDs {
			relations = append(relations, &TagRelation{
				TagID:      tagID,
				TargetID:   targetID,
				TargetType: targetType,
				Sort:       int32(i), // 按顺序排序
			})
		}

		if err := tx.Create(relations).Error; err != nil {
			return err
		}

		return nil
	})
}

// TagListFilter 标签筛选条件
type TagListFilter struct {
	Name     string
	Type     string
	Theme    string
	Page     int
	PageSize int
}

// ==================== Tag 定义相关实现 ====================

func (r *repository) GetTagByID(ctx context.Context, tagID uint64, tx ...*gorm.DB) (*Tag, error) {
	var tag Tag
	err := r.getDB(ctx, tx...).First(&tag, tagID).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *repository) GetTagByNameAndType(ctx context.Context, name, tagType string, tx ...*gorm.DB) (*Tag, error) {
	var tag Tag
	err := r.getDB(ctx, tx...).
		Where("name = ? AND type = ?", name, tagType).
		First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

func (r *repository) ListTags(ctx context.Context, filter TagListFilter, tx ...*gorm.DB) ([]*Tag, int64, error) {
	db := r.getDB(ctx, tx...).Model(&Tag{})

	if filter.Name != "" {
		db = db.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Type != "" {
		db = db.Where("type = ?", filter.Type)
	}
	if filter.Theme != "" {
		db = db.Where("theme = ?", filter.Theme)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit, offset := normalizePage(filter.Page, filter.PageSize)
	var list []*Tag
	err := db.Order("priority DESC, id ASC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *repository) UpdateTag(ctx context.Context, tagID uint64, values map[string]any, tx ...*gorm.DB) error {
	if tagID == 0 {
		return ErrTagNotFound
	}
	// 保护字段：不允许通过此接口修改 name 和 type，防止破坏唯一索引或业务逻辑
	delete(values, "name")
	delete(values, "type")
	delete(values, "id")

	res := r.getDB(ctx, tx...).Model(&Tag{}).Where("id = ?", tagID).Updates(values)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTagNotFound
	}
	return nil
}

func (r *repository) DeleteTag(ctx context.Context, tagID uint64, tx ...*gorm.DB) error {
	return r.withTx(ctx, func(innerTx *gorm.DB) error {
		// 1. 删除标签定义（软删除）
		if err := innerTx.Delete(&Tag{}, tagID).Error; err != nil {
			return err
		}
		// 2. 清理关联关系（物理删除关联，因为关联表通常不设软删）
		return innerTx.Where("tag_id = ?", tagID).Delete(&TagRelation{}).Error
	}, tx...)
}

// ==================== Tag 关联相关实现 ====================

func (r *repository) RemoveTagRelation(ctx context.Context, tagID uint64, targetType string, targetID uint64, tx ...*gorm.DB) error {
	res := r.getDB(ctx, tx...).
		Where("tag_id = ? AND target_type = ? AND target_id = ?", tagID, targetType, targetID).
		Delete(&TagRelation{})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("tag relation not found")
	}
	return nil
}

func (r *repository) ListTagRelationsByTarget(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) ([]*TagRelation, error) {
	var list []*TagRelation
	err := r.getDB(ctx, tx...).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("sort DESC, id ASC").
		Find(&list).Error
	return list, err
}
