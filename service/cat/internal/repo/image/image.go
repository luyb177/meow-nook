package image

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Image struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	TargetType  string `gorm:"type:varchar(32);not null;index:idx_target;comment:关联对象类型 cat_apply/cat_profile" json:"target_type"`
	TargetID    uint64 `gorm:"not null;index:idx_target;comment:关联对象ID" json:"target_id"`
	URL         string `gorm:"type:varchar(500);not null;comment:图片URL" json:"url"`
	Sort        int32  `gorm:"default:0;comment:排序" json:"sort"`
	IsCover     bool   `gorm:"default:false;comment:是否封面" json:"is_cover"`
	Description string `gorm:"type:varchar(255);comment:说明" json:"description"`
	UploaderID  uint64 `gorm:"default:0;comment:上传人ID" json:"uploader_id"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}
