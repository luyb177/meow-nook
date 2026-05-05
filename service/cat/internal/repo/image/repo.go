package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	TargetTypeCatApply   = "cat_apply"   // 申请
	TargetTypeCatProfile = "cat_profile" // 正式猫咪档案
)

var (
	ErrInvalidTarget = errors.New("invalid image target")
	ErrEmptyURL      = errors.New("image url is required")
	ErrMultiCover    = errors.New("only one cover image is allowed")
	ErrImageNotFound = errors.New("image not found")
)

type Repository interface {
	// Create 创建单张图片
	Create(ctx context.Context, img *Image, tx ...*gorm.DB) error

	// BatchCreate 批量创建图片
	BatchCreate(ctx context.Context, images []*Image, tx ...*gorm.DB) error

	// ReplaceByTarget 批量替换某个对象的图片：先软删原有图片，再创建新图片
	ReplaceByTarget(ctx context.Context, targetType string, targetID uint64, images []*Image, tx ...*gorm.DB) error

	// DeleteByTarget 删除某个对象下的全部图片
	DeleteByTarget(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) error

	// ListByTarget 查询某个对象下的图片
	ListByTarget(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) ([]*Image, error)

	// GetCover 获取封面图，如果没有封面，则返回排序最高的第一张
	GetCover(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) (*Image, error)

	// SetCover 设置封面图
	SetCover(ctx context.Context, targetType string, targetID uint64, imageID uint64, tx ...*gorm.DB) error

	// CopyByTarget 复制图片归属，例如 cat_apply -> cat_profile
	CopyByTarget(ctx context.Context, fromType string, fromID uint64, toType string, toID uint64, uploaderID uint64, tx ...*gorm.DB) error

	// MoveByTarget 移动图片归属，例如 cat_apply -> cat_profile
	// 如果想保留申请图片快照，建议用 CopyByTarget
	MoveByTarget(ctx context.Context, fromType string, fromID uint64, toType string, toID uint64, tx ...*gorm.DB) error
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
func (r *repository) Create(ctx context.Context, img *Image, tx ...*gorm.DB) error {
	if img == nil {
		return nil // 或者返回 errors.New("image object is nil")
	}
	if err := r.validateImage(img); err != nil {
		return err
	}
	return r.getDB(ctx, tx...).Create(img).Error
}

func (r *repository) BatchCreate(ctx context.Context, images []*Image, tx ...*gorm.DB) error {
	if len(images) == 0 {
		return nil
	}
	if err := r.validateImages(images); err != nil {
		return err
	}
	return r.getDB(ctx, tx...).Create(&images).Error
}

func (r *repository) ReplaceByTarget(ctx context.Context, targetType string, targetID uint64, images []*Image, tx ...*gorm.DB) error {
	if err := r.validateTarget(targetType, targetID); err != nil {
		return err
	}

	return r.withTx(ctx, func(innerTx *gorm.DB) error {
		// 1. 软删除原有图片
		if err := innerTx.
			Where("target_type = ? AND target_id = ?", targetType, targetID).
			Delete(&Image{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing images: %w", err)
		}

		// 2. 如果没有新图片，则直接返回
		if len(images) == 0 {
			return nil
		}

		// 3. 重新设置新图片的 TargetType 和 TargetID
		for _, img := range images {
			img.ID = 0 // 确保是新记录
			img.TargetType = targetType
			img.TargetID = targetID
		}

		// 4. 校验新图片
		if err := r.validateImages(images); err != nil {
			return err
		}

		// 5. 批量创建新图片
		if err := innerTx.Create(&images).Error; err != nil {
			return fmt.Errorf("failed to create new images: %w", err)
		}
		return nil
	}, tx...)
}

func (r *repository) DeleteByTarget(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) error {
	if err := r.validateTarget(targetType, targetID); err != nil {
		return err
	}
	return r.getDB(ctx, tx...).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Delete(&Image{}).Error
}

func (r *repository) ListByTarget(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) ([]*Image, error) {
	if err := r.validateTarget(targetType, targetID); err != nil {
		return nil, err
	}

	var list []*Image
	err := r.getDB(ctx, tx...).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("is_cover DESC, sort DESC, id ASC"). // 封面优先，然后排序，然后ID
		Find(&list).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []*Image{}, nil // 没有找到返回空列表，而不是错误
	}
	return list, err
}
func (r *repository) GetCover(ctx context.Context, targetType string, targetID uint64, tx ...*gorm.DB) (*Image, error) {
	if err := r.validateTarget(targetType, targetID); err != nil {
		return nil, err
	}

	var img Image
	db := r.getDB(ctx, tx...)

	// 优先查询设置为封面的图片
	err := db.Where("target_type = ? AND target_id = ? AND is_cover = ?", targetType, targetID, true).First(&img).Error
	if err == nil {
		return &img, nil
	}

	// 如果没有找到封面图片，或者找到了但是有其他错误
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // 其他数据库错误
	}

	// 如果没有明确的封面，就按排序查找第一张
	err = db.Where("target_type = ? AND target_id = ?", targetType, targetID).Order("sort DESC, id ASC").First(&img).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 没有找到任何图片
	}
	return &img, err
}

func (r *repository) SetCover(ctx context.Context, targetType string, targetID uint64, imageID uint64, tx ...*gorm.DB) error {
	if err := r.validateTarget(targetType, targetID); err != nil {
		return err
	}
	if imageID == 0 {
		return ErrImageNotFound // 尝试设置 ID 为 0 的图片为封面
	}

	return r.withTx(ctx, func(innerTx *gorm.DB) error {
		// 1. 取消该 Target 下所有图片的封面状态
		if err := innerTx.Model(&Image{}).
			Where("target_type = ? AND target_id = ?", targetType, targetID).
			Update("is_cover", false).Error; err != nil {
			return fmt.Errorf("failed to clear existing cover: %w", err)
		}

		// 2. 设置指定图片为封面
		result := innerTx.Model(&Image{}).
			Where("id = ? AND target_type = ? AND target_id = ?", imageID, targetType, targetID).
			Update("is_cover", true)

		if result.Error != nil {
			return fmt.Errorf("failed to set new cover: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrImageNotFound // 指定的图片不存在或不属于该 target
		}
		return nil
	}, tx...)
}

func (r *repository) CopyByTarget(ctx context.Context, fromType string, fromID uint64, toType string, toID uint64, uploaderID uint64, tx ...*gorm.DB) error {
	if err := r.validateTarget(fromType, fromID); err != nil {
		return err
	}
	if err := r.validateTarget(toType, toID); err != nil {
		return err
	}

	return r.withTx(ctx, func(innerTx *gorm.DB) error {
		// 1. 查询源 Target 的所有图片
		sourceImages, err := r.ListByTarget(ctx, fromType, fromID, innerTx)
		if err != nil {
			return fmt.Errorf("failed to list source images: %w", err)
		}
		if len(sourceImages) == 0 {
			return nil // 没有源图片可复制
		}

		// 2. 构建新图片列表
		newImages := make([]*Image, 0, len(sourceImages))
		for _, oldImg := range sourceImages {
			newImages = append(newImages, &Image{
				TargetType:  toType,
				TargetID:    toID,
				URL:         oldImg.URL,
				Sort:        oldImg.Sort,
				IsCover:     oldImg.IsCover,
				Description: oldImg.Description,
				UploaderID:  uploaderID, // 使用传入的 uploaderID
			})
		}

		// 3. 批量创建新图片
		if err := innerTx.Create(&newImages).Error; err != nil {
			return fmt.Errorf("failed to create copied images: %w", err)
		}
		return nil
	}, tx...)
}
func (r *repository) MoveByTarget(ctx context.Context, fromType string, fromID uint64, toType string, toID uint64, tx ...*gorm.DB) error {
	if err := r.validateTarget(fromType, fromID); err != nil {
		return err
	}
	if err := r.validateTarget(toType, toID); err != nil {
		return err
	}

	return r.getDB(ctx, tx...).Model(&Image{}).
		Where("target_type = ? AND target_id = ?", fromType, fromID).
		Updates(map[string]interface{}{
			"target_type": toType,
			"target_id":   toID,
		}).Error
}
