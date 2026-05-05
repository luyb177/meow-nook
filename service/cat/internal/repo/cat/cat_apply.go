package cat

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/soft_delete"
)

// ==================== 状态常量 ====================

const (
	ApplyStatusPending  = "pending"
	ApplyStatusApproved = "approved"
	ApplyStatusRejected = "rejected"
	ApplyStatusCanceled = "canceled"
)

// ==================== 错误定义 ====================

var (
	ErrApplyNotPending         = errors.New("apply is not pending")
	ErrInvalidApplyStatus      = errors.New("invalid apply status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrApplicantUserIDRequired = errors.New("applicant_user_id is required")
	ErrApplyNameRequired       = errors.New("apply name is required")
)

// ==================== 过滤条件 ====================

type ApplyListFilter struct {
	Status          *string // 可选
	ApplicantUserID uint64  // 可选
	ReviewerID      uint64  // 可选
	Keyword         string  // 可选，按名称模糊搜索

	Page     int
	PageSize int
}

type CatCreateApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`

	// 审核通过后关联的正式档案
	CatID uint64 `gorm:"index;default:0;comment:审核通过后生成的猫咪档案ID" json:"cat_id"`

	// 基础信息
	Name        string `gorm:"type:varchar(64);not null;comment:猫咪名称" json:"name"`
	Gender      string `gorm:"type:varchar(20);default:'unknown';comment:性别" json:"gender"`
	BodySize    string `gorm:"type:varchar(20);default:'medium';comment:体型" json:"body_size"`
	AgeStage    string `gorm:"type:varchar(30);default:'young';comment:年龄阶段" json:"age_stage"`
	Description string `gorm:"type:text;comment:简介" json:"description"`

	// 位置信息
	DiscoveryAddress string  `gorm:"type:varchar(255);comment:发现地址" json:"discovery_address"`
	Longitude        float64 `gorm:"type:decimal(10,6);default:0" json:"longitude"`
	Latitude         float64 `gorm:"type:decimal(10,6);default:0" json:"latitude"`

	// 申请人信息
	ApplicantUserID uint64 `gorm:"index;not null;comment:申请人ID" json:"applicant_user_id"`

	// 审核信息
	Status       string     `gorm:"type:varchar(30);default:'pending';index;comment:pending/approved/rejected/canceled" json:"status"`
	RejectReason string     `gorm:"type:text;comment:驳回原因" json:"reject_reason"`
	ReviewerID   uint64     `gorm:"default:0;comment:审核人ID" json:"reviewer_id"`
	ReviewedAt   *time.Time `gorm:"comment:审核时间" json:"reviewed_at,omitempty"`

	// 取消信息
	CanceledAt   *time.Time `gorm:"comment:取消时间" json:"canceled_at,omitempty"`
	CancelReason string     `gorm:"type:text;comment:取消原因" json:"cancel_reason"`

	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatCreateApply) TableName() string {
	return "cat_create_applies"
}

// ==================== 申请相关实现 ====================

func (r *repository) CreateApply(ctx context.Context, apply *CatCreateApply, tx ...*gorm.DB) error {
	if apply == nil {
		return errors.New("apply is nil")
	}
	if apply.Name == "" {
		return ErrApplyNameRequired
	}
	if apply.ApplicantUserID == 0 {
		return ErrApplicantUserIDRequired
	}
	if apply.Status == "" {
		apply.Status = ApplyStatusPending
	}

	return r.getDB(ctx, tx...).Create(apply).Error
}

func (r *repository) GetApplyByID(ctx context.Context, applyID uint64, tx ...*gorm.DB) (*CatCreateApply, error) {
	var apply CatCreateApply
	err := r.getDB(ctx, tx...).
		Where("id = ?", applyID).
		First(&apply).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repository) GetApplyByIDForUpdate(ctx context.Context, applyID uint64, tx ...*gorm.DB) (*CatCreateApply, error) {
	var apply CatCreateApply
	err := r.getDB(ctx, tx...).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", applyID).
		First(&apply).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repository) ListApplies(ctx context.Context, filter ApplyListFilter, tx ...*gorm.DB) ([]*CatCreateApply, int64, error) {
	db := r.getDB(ctx, tx...).Model(&CatCreateApply{})

	if filter.Status != nil && *filter.Status != "" {
		db = db.Where("status = ?", *filter.Status)
	}
	if filter.ApplicantUserID > 0 {
		db = db.Where("applicant_user_id = ?", filter.ApplicantUserID)
	}
	if filter.ReviewerID > 0 {
		db = db.Where("reviewer_id = ?", filter.ReviewerID)
	}
	if filter.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+filter.Keyword+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit, offset := normalizePage(filter.Page, filter.PageSize)

	var list []*CatCreateApply
	err := db.Order("id DESC").Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) UpdateApply(ctx context.Context, applyID uint64, values map[string]any, tx ...*gorm.DB) error {
	if applyID == 0 {
		return gorm.ErrRecordNotFound
	}
	if len(values) == 0 {
		return nil
	}

	res := r.getDB(ctx, tx...).
		Model(&CatCreateApply{}).
		Where("id = ? AND status = ?", applyID, ApplyStatusPending).
		Updates(values)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (r *repository) CancelApply(ctx context.Context, applyID uint64, applicantUserID uint64, reason string, tx ...*gorm.DB) error {
	if applicantUserID == 0 {
		return ErrApplicantUserIDRequired
	}

	now := time.Now()
	res := r.getDB(ctx, tx...).
		Model(&CatCreateApply{}).
		Where("id = ? AND status = ? AND applicant_user_id = ?", applyID, ApplyStatusPending, applicantUserID).
		Updates(map[string]any{
			"status":        ApplyStatusCanceled,
			"cancel_reason": reason,
			"canceled_at":   &now,
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (r *repository) RejectApply(ctx context.Context, applyID uint64, reviewerID uint64, reason string, tx ...*gorm.DB) error {
	now := time.Now()
	res := r.getDB(ctx, tx...).
		Model(&CatCreateApply{}).
		Where("id = ? AND status = ?", applyID, ApplyStatusPending).
		Updates(map[string]any{
			"status":        ApplyStatusRejected,
			"reviewer_id":   reviewerID,
			"reject_reason": reason,
			"reviewed_at":   &now,
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (r *repository) ApproveApply(ctx context.Context, applyID uint64, reviewerID uint64, catID uint64, tx ...*gorm.DB) error {
	if catID == 0 {
		return ErrCatCodeRequired
	}

	now := time.Now()
	res := r.getDB(ctx, tx...).
		Model(&CatCreateApply{}).
		Where("id = ? AND status = ?", applyID, ApplyStatusPending).
		Updates(map[string]any{
			"status":        ApplyStatusApproved,
			"cat_id":        catID,
			"reviewer_id":   reviewerID,
			"reviewed_at":   &now,
			"reject_reason": "",
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (r *repository) DeleteApply(ctx context.Context, applyID uint64, tx ...*gorm.DB) error {
	res := r.getDB(ctx, tx...).
		Where("id = ?", applyID).
		Delete(&CatCreateApply{})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
