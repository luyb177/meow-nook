package adoption

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/soft_delete"
)

// ==================== 常量 ====================

const (
	AdoptionStatusActive   = "active"   // 正常领养中
	AdoptionStatusReturned = "returned" // 已退回平台
	AdoptionStatusExpired  = "expired"  // 已退养（超期未回访等）
)

// ==================== 错误定义 ====================

var (
	ErrAdoptionNotFound      = errors.New("领养记录不存在")
	ErrAdoptionAlreadyExists = errors.New("该猫咪已有领养记录")
	ErrVisitAlreadyCompleted = errors.New("该次回访已完成")
	ErrVisitNotDue           = errors.New("未到回访时间")
	ErrReturnReasonRequired  = errors.New("退回原因不能为空")
)

// ==================== 数据模型 ====================

type Adoption struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// ==================== 关联信息 ====================
	CatID     uint64 `gorm:"not null;uniqueIndex;comment:猫咪ID（一只猫只有一条领养记录）" json:"cat_id"`
	AdopterID uint64 `gorm:"not null;index;comment:领养人ID" json:"adopter_id"`
	ApplyID   uint64 `gorm:"default:0;index;comment:来源申请ID，0表示管理员直接领养" json:"apply_id"`
	CreatorID uint64 `gorm:"not null;comment:创建人ID（管理员）" json:"creator_id"`

	// ==================== 领养基础信息 ====================
	Status      string     `gorm:"type:varchar(30);default:'active';index;comment:领养状态 active/returned/expired" json:"status"`
	AgreedAt    *time.Time `gorm:"comment:协议签订时间" json:"agreed_at"`
	AgreementNo string     `gorm:"type:varchar(64);comment:线下协议编号" json:"agreement_no"`
	AdoptedAt   *time.Time `gorm:"comment:正式领养时间" json:"adopted_at"`

	// ==================== 家访信息 ====================
	// 第一次家访（领养前/领养后1周内）
	HomeVisitAt     *time.Time `gorm:"comment:首次家访时间" json:"home_visit_at"`
	HomeVisitUserID uint64     `gorm:"default:0;comment:家访志愿者ID" json:"home_visit_user_id"`
	HomeVisitPhotos string     `gorm:"type:text;comment:家访照片(JSON数组)" json:"home_visit_photos"`
	HomeVisitRemark string     `gorm:"type:text;comment:家访备注" json:"home_visit_remark"`

	// ==================== 回访记录 ====================
	// 1周回访
	VisitOneWeekAt     *time.Time `gorm:"comment:1周回访时间" json:"visit_one_week_at"`
	VisitOneWeekUserID uint64     `gorm:"default:0;comment:1周回访志愿者ID" json:"visit_one_week_user_id"`
	VisitOneWeekPhotos string     `gorm:"type:text;comment:1周回访照片" json:"visit_one_week_photos"`
	VisitOneWeekRemark string     `gorm:"type:text;comment:1周回访备注" json:"visit_one_week_remark"`

	// 1月回访
	VisitOneMonthAt     *time.Time `gorm:"comment:1月回访时间" json:"visit_one_month_at"`
	VisitOneMonthUserID uint64     `gorm:"default:0;comment:1月回访志愿者ID" json:"visit_one_month_user_id"`
	VisitOneMonthPhotos string     `gorm:"type:text;comment:1月回访照片" json:"visit_one_month_photos"`
	VisitOneMonthRemark string     `gorm:"type:text;comment:1月回访备注" json:"visit_one_month_remark"`

	// 3月回访
	VisitThreeMonthAt     *time.Time `gorm:"comment:3月回访时间" json:"visit_three_month_at"`
	VisitThreeMonthUserID uint64     `gorm:"default:0;comment:3月回访志愿者ID" json:"visit_three_month_user_id"`
	VisitThreeMonthPhotos string     `gorm:"type:text;comment:3月回访照片" json:"visit_three_month_photos"`
	VisitThreeMonthRemark string     `gorm:"type:text;comment:3月回访备注" json:"visit_three_month_remark"`

	// 6月回访
	VisitSixMonthAt     *time.Time `gorm:"comment:6月回访时间" json:"visit_six_month_at"`
	VisitSixMonthUserID uint64     `gorm:"default:0;comment:6月回访志愿者ID" json:"visit_six_month_user_id"`
	VisitSixMonthPhotos string     `gorm:"type:text;comment:6月回访照片" json:"visit_six_month_photos"`
	VisitSixMonthRemark string     `gorm:"type:text;comment:6月回访备注" json:"visit_six_month_remark"`

	// ==================== 退回/退养信息 ====================
	IsReturned       bool       `gorm:"default:false;comment:是否已退回" json:"is_returned"`
	ReturnedAt       *time.Time `gorm:"comment:退回时间" json:"returned_at"`
	ReturnReason     string     `gorm:"type:text;comment:退回原因" json:"return_reason"`
	ReturnedToUserID uint64     `gorm:"default:0;comment:退回给谁，0表示退回平台" json:"returned_to_user_id"`
	ReturnPhotos     string     `gorm:"type:text;comment:退回时照片" json:"return_photos"`

	// ==================== 备注 ====================
	Note string `gorm:"type:text;comment:备注" json:"note"`

	// ==================== 通用字段 ====================
	CreatedAt time.Time             `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (Adoption) TableName() string {
	return "adoptions"
}

// ==================== 查询过滤器 ====================

type AdoptionFilter struct {
	CatID     uint64 // 按猫咪筛选
	AdopterID uint64 // 按领养人筛选
	Status    string // active/returned/expired

	// 回访状态筛选
	VisitOneWeekDone    *bool
	VisitOneMonthDone   *bool
	VisitThreeMonthDone *bool
	VisitSixMonthDone   *bool

	CreatedAtStart *time.Time
	CreatedAtEnd   *time.Time

	Page     int
	PageSize int
}

// ==================== Adoption 方法实现 ====================

func (r *repository) CreateAdoption(ctx context.Context, adoption *Adoption, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		// 检查猫咪是否已有有效领养记录
		var existCount int64
		if err := tx.Model(&Adoption{}).Where("cat_id = ? AND status != ?", adoption.CatID, AdoptionStatusReturned).Count(&existCount).Error; err != nil {
			return err
		}
		if existCount > 0 {
			return ErrAdoptionAlreadyExists
		}

		return tx.Create(adoption).Error
	})
}

func (r *repository) GetAdoptionByID(ctx context.Context, adoptionID uint64, txs ...*gorm.DB) (*Adoption, error) {
	tx := r.getDB(ctx, txs...)
	var adoption Adoption
	err := tx.First(&adoption, adoptionID).Error
	if err != nil {
		return nil, err
	}
	return &adoption, nil
}

func (r *repository) GetAdoptionByCatID(ctx context.Context, catID uint64, txs ...*gorm.DB) (*Adoption, error) {
	tx := r.getDB(ctx, txs...)
	var adoption Adoption

	// 查找非退回状态的记录
	err := tx.Where("cat_id = ? AND status != ?", catID, AdoptionStatusReturned).Order("created_at DESC").First(&adoption).Error
	if err != nil {
		return nil, err
	}
	return &adoption, nil
}

func (r *repository) ListAdoptions(ctx context.Context, filter *AdoptionFilter, txs ...*gorm.DB) ([]*Adoption, int64, error) {
	tx := r.getDB(ctx, txs...)
	var total int64

	query := tx.Model(&Adoption{})

	if filter.CatID > 0 {
		query = query.Where("cat_id = ?", filter.CatID)
	}
	if filter.AdopterID > 0 {
		query = query.Where("adopter_id = ?", filter.AdopterID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.VisitOneWeekDone != nil && *filter.VisitOneWeekDone {
		query = query.Where("visit_one_week_at IS NOT NULL")
	} else if filter.VisitOneWeekDone != nil && !*filter.VisitOneWeekDone {
		query = query.Where("visit_one_week_at IS NULL")
	}
	if filter.VisitOneMonthDone != nil && *filter.VisitOneMonthDone {
		query = query.Where("visit_one_month_at IS NOT NULL")
	} else if filter.VisitOneMonthDone != nil && !*filter.VisitOneMonthDone {
		query = query.Where("visit_one_month_at IS NULL")
	}
	if filter.VisitThreeMonthDone != nil && *filter.VisitThreeMonthDone {
		query = query.Where("visit_three_month_at IS NOT NULL")
	} else if filter.VisitThreeMonthDone != nil && !*filter.VisitThreeMonthDone {
		query = query.Where("visit_three_month_at IS NULL")
	}
	if filter.VisitSixMonthDone != nil && *filter.VisitSixMonthDone {
		query = query.Where("visit_six_month_at IS NOT NULL")
	} else if filter.VisitSixMonthDone != nil && !*filter.VisitSixMonthDone {
		query = query.Where("visit_six_month_at IS NULL")
	}
	if filter.CreatedAtStart != nil {
		query = query.Where("created_at >= ?", *filter.CreatedAtStart)
	}
	if filter.CreatedAtEnd != nil {
		query = query.Where("created_at <= ?", *filter.CreatedAtEnd)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (filter.Page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	var adoptions []*Adoption
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&adoptions).Error
	return adoptions, total, err
}

func (r *repository) UpdateAdoption(ctx context.Context, adoptionID uint64, values map[string]any, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)
	return tx.Model(&Adoption{}).Where("id = ?", adoptionID).Updates(values).Error
}

func (r *repository) RecordHomeVisit(ctx context.Context, adoptionID uint64, homeVisitAt time.Time, visitorID uint64, photos, remark string, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		var adoption Adoption
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&adoption, adoptionID).Error; err != nil {
			return err
		}

		// 检查家访是否已完成
		if adoption.HomeVisitAt != nil && !adoption.HomeVisitAt.IsZero() {
			return ErrVisitAlreadyCompleted
		}

		adoption.HomeVisitAt = &homeVisitAt
		adoption.HomeVisitUserID = visitorID
		adoption.HomeVisitPhotos = photos
		adoption.HomeVisitRemark = remark

		return tx.Save(&adoption).Error
	})
}

func (r *repository) RecordFollowUpVisit(ctx context.Context, adoptionID uint64, visitType int, visitorID uint64, visitedAt time.Time, photos, remark string, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		var adoption Adoption
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&adoption, adoptionID).Error; err != nil {
			return err
		}

		data := make(map[string]interface{})

		switch visitType {
		case 1: // 1 周
			if adoption.VisitOneWeekAt != nil && !adoption.VisitOneWeekAt.IsZero() {
				return ErrVisitAlreadyCompleted
			}
			data["visit_one_week_at"] = visitedAt
			data["visit_one_week_user_id"] = visitorID
			data["visit_one_week_photos"] = photos
			data["visit_one_week_remark"] = remark
		case 2: // 1 月
			if adoption.VisitOneMonthAt != nil && !adoption.VisitOneMonthAt.IsZero() {
				return ErrVisitAlreadyCompleted
			}
			data["visit_one_month_at"] = visitedAt
			data["visit_one_month_user_id"] = visitorID
			data["visit_one_month_photos"] = photos
			data["visit_one_month_remark"] = remark
		case 3: // 3 月
			if adoption.VisitThreeMonthAt != nil && !adoption.VisitThreeMonthAt.IsZero() {
				return ErrVisitAlreadyCompleted
			}
			data["visit_three_month_at"] = visitedAt
			data["visit_three_month_user_id"] = visitorID
			data["visit_three_month_photos"] = photos
			data["visit_three_month_remark"] = remark
		case 4: // 6 月
			if adoption.VisitSixMonthAt != nil && !adoption.VisitSixMonthAt.IsZero() {
				return ErrVisitAlreadyCompleted
			}
			data["visit_six_month_at"] = visitedAt
			data["visit_six_month_user_id"] = visitorID
			data["visit_six_month_photos"] = photos
			data["visit_six_month_remark"] = remark
		default:
			return errors.New("invalid visit type")
		}

		return tx.Model(&adoption).Updates(data).Error
	})
}

func (r *repository) UpdateReturnStatus(ctx context.Context, adoptionID uint64, returned bool, returnedToUserID uint64, returnReason, photos string, returnedAt time.Time, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		var adoption Adoption
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&adoption, adoptionID).Error; err != nil {
			return err
		}

		// 如果标记为退回，校验原因
		if returned && returnReason == "" {
			return ErrReturnReasonRequired
		}

		adoption.IsReturned = returned
		adoption.ReturnedToUserID = returnedToUserID
		adoption.ReturnReason = returnReason
		adoption.ReturnPhotos = photos

		if returned {
			adoption.ReturnedAt = &returnedAt
			adoption.Status = AdoptionStatusReturned
		} else {
			adoption.Status = AdoptionStatusActive
		}

		return tx.Save(&adoption).Error
	})
}

func (r *repository) DeleteAdoption(ctx context.Context, adoptionID uint64, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)
	err := tx.Delete(&Adoption{}, adoptionID).Error
	return err
}
