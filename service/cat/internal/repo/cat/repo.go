package cat

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Repository defines all data access operations for the cat domain.
type Repository interface {
	// Cat
	CreateCat(ctx context.Context, c *Cat) error
	GetCatByID(ctx context.Context, id uint64) (*Cat, error)
	UpdateCat(ctx context.Context, id uint64, updates map[string]interface{}) error
	ListCats(ctx context.Context, lastID uint64, pageSize int) ([]*Cat, bool, error)

	// CatImage
	CreateCatImages(ctx context.Context, images []*CatImage) error
	ListCatImages(ctx context.Context, catID uint64) ([]*CatImage, error)

	// CatCreateApply
	CreateCatCreateApply(ctx context.Context, apply *CatCreateApply) error
	GetCatCreateApplyByID(ctx context.Context, id uint64) (*CatCreateApply, error)
	GetCatCreateApplyByIDAndUser(ctx context.Context, id, userID uint64) (*CatCreateApply, error)
	ListCatCreateApplies(ctx context.Context, lastID uint64, pageSize int) ([]*CatCreateApply, bool, error)
	UpdateCatCreateApplyStatus(ctx context.Context, id uint64, status, rejectReason string, reviewerID uint64) error

	// CatUpdateApply
	CreateCatUpdateApply(ctx context.Context, apply *CatUpdateApply) error
	GetCatUpdateApplyByID(ctx context.Context, id uint64) (*CatUpdateApply, error)
	GetCatUpdateApplyByIDAndUser(ctx context.Context, id, userID uint64) (*CatUpdateApply, error)
	ListCatUpdateApplies(ctx context.Context, lastID uint64, pageSize int) ([]*CatUpdateApply, bool, error)
	UpdateCatUpdateApplyStatus(ctx context.Context, id uint64, status, rejectReason string, reviewerID uint64) error

	// CatMedicalApply
	CreateCatMedicalApply(ctx context.Context, apply *CatMedicalApply) error
	GetCatMedicalApplyByID(ctx context.Context, id uint64) (*CatMedicalApply, error)
	UpdateCatMedicalApplyStatus(ctx context.Context, id uint64, status, rejectReason string, reviewerID uint64) error

	// CatMedicalRecord
	CreateCatMedicalRecord(ctx context.Context, record *CatMedicalRecord) error

	// CatRescueApply
	CreateCatRescueApply(ctx context.Context, apply *CatRescueApply) error

	// CatTaskApply
	CreateCatTaskApply(ctx context.Context, apply *CatTaskApply) error
	GetCatTaskApplyByID(ctx context.Context, id uint64) (*CatTaskApply, error)
	ListCatTaskApplies(ctx context.Context, lastID uint64, pageSize int) ([]*CatTaskApply, bool, error)
	UpdateCatTaskApply(ctx context.Context, id uint64, updates map[string]interface{}) error

	// CatTask
	CreateCatTask(ctx context.Context, task *CatTask) error
	GetCatTaskByID(ctx context.Context, id uint64) (*CatTask, error)
	ListCatTasks(ctx context.Context, lastID uint64, pageSize int) ([]*CatTask, bool, error)
	UpdateCatTask(ctx context.Context, id uint64, updates map[string]interface{}) error

	// CatTaskClaim
	CreateCatTaskClaim(ctx context.Context, claim *CatTaskClaim) error
	GetActiveClaimByTaskAndUser(ctx context.Context, taskID, userID uint64) (*CatTaskClaim, error)
	UpdateCatTaskClaim(ctx context.Context, id uint64, updates map[string]interface{}) error

	// CatAdoption
	UpsertCatAdoption(ctx context.Context, catID uint64, status string) error

	// CatTaskProgress
	CreateCatTaskProgress(ctx context.Context, progress *CatTaskProgress) error
}

type repo struct {
	db     *gorm.DB
	client *redis.Client
}

func NewRepository(db *gorm.DB, client *redis.Client) Repository {
	return &repo{db: db, client: client}
}

// ── Cat ──────────────────────────────────────────────────────────────────────

func (r *repo) CreateCat(ctx context.Context, c *Cat) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *repo) GetCatByID(ctx context.Context, id uint64) (*Cat, error) {
	var c Cat
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repo) UpdateCat(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&Cat{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repo) ListCats(ctx context.Context, lastID uint64, pageSize int) ([]*Cat, bool, error) {
	limit := pageSize + 1
	var cats []*Cat
	q := r.db.WithContext(ctx).Order("id DESC").Limit(limit)
	if lastID > 0 {
		q = q.Where("id < ?", lastID)
	}
	if err := q.Find(&cats).Error; err != nil {
		return nil, false, err
	}
	hasMore := len(cats) == limit
	if hasMore {
		cats = cats[:pageSize]
	}
	return cats, hasMore, nil
}

// ── CatImage ─────────────────────────────────────────────────────────────────

func (r *repo) CreateCatImages(ctx context.Context, images []*CatImage) error {
	if len(images) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&images).Error
}

func (r *repo) ListCatImages(ctx context.Context, catID uint64) ([]*CatImage, error) {
	var images []*CatImage
	if err := r.db.WithContext(ctx).Where("cat_id = ?", catID).Order("sort ASC").Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

// ── CatCreateApply ────────────────────────────────────────────────────────────

func (r *repo) CreateCatCreateApply(ctx context.Context, apply *CatCreateApply) error {
	return r.db.WithContext(ctx).Create(apply).Error
}

func (r *repo) GetCatCreateApplyByID(ctx context.Context, id uint64) (*CatCreateApply, error) {
	var apply CatCreateApply
	if err := r.db.WithContext(ctx).First(&apply, id).Error; err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repo) GetCatCreateApplyByIDAndUser(ctx context.Context, id, userID uint64) (*CatCreateApply, error) {
	var apply CatCreateApply
	if err := r.db.WithContext(ctx).Where("id = ? AND applicant_user_id = ?", id, userID).First(&apply).Error; err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repo) ListCatCreateApplies(ctx context.Context, lastID uint64, pageSize int) ([]*CatCreateApply, bool, error) {
	limit := pageSize + 1
	var applies []*CatCreateApply
	q := r.db.WithContext(ctx).Order("id DESC").Limit(limit)
	if lastID > 0 {
		q = q.Where("id < ?", lastID)
	}
	if err := q.Find(&applies).Error; err != nil {
		return nil, false, err
	}
	hasMore := len(applies) == limit
	if hasMore {
		applies = applies[:pageSize]
	}
	return applies, hasMore, nil
}

func (r *repo) UpdateCatCreateApplyStatus(ctx context.Context, id uint64, status, rejectReason string, reviewerID uint64) error {
	return r.db.WithContext(ctx).Model(&CatCreateApply{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"reject_reason": rejectReason,
		"reviewer_id":   reviewerID,
	}).Error
}

// ── CatUpdateApply ────────────────────────────────────────────────────────────

func (r *repo) CreateCatUpdateApply(ctx context.Context, apply *CatUpdateApply) error {
	return r.db.WithContext(ctx).Create(apply).Error
}

func (r *repo) GetCatUpdateApplyByID(ctx context.Context, id uint64) (*CatUpdateApply, error) {
	var apply CatUpdateApply
	if err := r.db.WithContext(ctx).First(&apply, id).Error; err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repo) GetCatUpdateApplyByIDAndUser(ctx context.Context, id, userID uint64) (*CatUpdateApply, error) {
	var apply CatUpdateApply
	if err := r.db.WithContext(ctx).Where("id = ? AND applicant_user_id = ?", id, userID).First(&apply).Error; err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repo) ListCatUpdateApplies(ctx context.Context, lastID uint64, pageSize int) ([]*CatUpdateApply, bool, error) {
	limit := pageSize + 1
	var applies []*CatUpdateApply
	q := r.db.WithContext(ctx).Order("id DESC").Limit(limit)
	if lastID > 0 {
		q = q.Where("id < ?", lastID)
	}
	if err := q.Find(&applies).Error; err != nil {
		return nil, false, err
	}
	hasMore := len(applies) == limit
	if hasMore {
		applies = applies[:pageSize]
	}
	return applies, hasMore, nil
}

func (r *repo) UpdateCatUpdateApplyStatus(ctx context.Context, id uint64, status, rejectReason string, reviewerID uint64) error {
	return r.db.WithContext(ctx).Model(&CatUpdateApply{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"reject_reason": rejectReason,
		"reviewer_id":   reviewerID,
	}).Error
}

// ── CatMedicalApply ───────────────────────────────────────────────────────────

func (r *repo) CreateCatMedicalApply(ctx context.Context, apply *CatMedicalApply) error {
	return r.db.WithContext(ctx).Create(apply).Error
}

func (r *repo) GetCatMedicalApplyByID(ctx context.Context, id uint64) (*CatMedicalApply, error) {
	var apply CatMedicalApply
	if err := r.db.WithContext(ctx).First(&apply, id).Error; err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repo) UpdateCatMedicalApplyStatus(ctx context.Context, id uint64, status, rejectReason string, reviewerID uint64) error {
	return r.db.WithContext(ctx).Model(&CatMedicalApply{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"reject_reason": rejectReason,
		"reviewer_id":   reviewerID,
	}).Error
}

// ── CatMedicalRecord ──────────────────────────────────────────────────────────

func (r *repo) CreateCatMedicalRecord(ctx context.Context, record *CatMedicalRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// ── CatRescueApply ────────────────────────────────────────────────────────────

func (r *repo) CreateCatRescueApply(ctx context.Context, apply *CatRescueApply) error {
	return r.db.WithContext(ctx).Create(apply).Error
}

// ── CatTaskApply ──────────────────────────────────────────────────────────────

func (r *repo) CreateCatTaskApply(ctx context.Context, apply *CatTaskApply) error {
	return r.db.WithContext(ctx).Create(apply).Error
}

func (r *repo) GetCatTaskApplyByID(ctx context.Context, id uint64) (*CatTaskApply, error) {
	var apply CatTaskApply
	if err := r.db.WithContext(ctx).First(&apply, id).Error; err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repo) ListCatTaskApplies(ctx context.Context, lastID uint64, pageSize int) ([]*CatTaskApply, bool, error) {
	limit := pageSize + 1
	var applies []*CatTaskApply
	q := r.db.WithContext(ctx).Order("id DESC").Limit(limit)
	if lastID > 0 {
		q = q.Where("id < ?", lastID)
	}
	if err := q.Find(&applies).Error; err != nil {
		return nil, false, err
	}
	hasMore := len(applies) == limit
	if hasMore {
		applies = applies[:pageSize]
	}
	return applies, hasMore, nil
}

func (r *repo) UpdateCatTaskApply(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&CatTaskApply{}).Where("id = ?", id).Updates(updates).Error
}

// ── CatTask ───────────────────────────────────────────────────────────────────

func (r *repo) CreateCatTask(ctx context.Context, task *CatTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *repo) GetCatTaskByID(ctx context.Context, id uint64) (*CatTask, error) {
	var task CatTask
	if err := r.db.WithContext(ctx).First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *repo) ListCatTasks(ctx context.Context, lastID uint64, pageSize int) ([]*CatTask, bool, error) {
	limit := pageSize + 1
	var tasks []*CatTask
	q := r.db.WithContext(ctx).Order("id DESC").Limit(limit)
	if lastID > 0 {
		q = q.Where("id < ?", lastID)
	}
	if err := q.Find(&tasks).Error; err != nil {
		return nil, false, err
	}
	hasMore := len(tasks) == limit
	if hasMore {
		tasks = tasks[:pageSize]
	}
	return tasks, hasMore, nil
}

func (r *repo) UpdateCatTask(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&CatTask{}).Where("id = ?", id).Updates(updates).Error
}

// ── CatTaskClaim ──────────────────────────────────────────────────────────────

func (r *repo) CreateCatTaskClaim(ctx context.Context, claim *CatTaskClaim) error {
	return r.db.WithContext(ctx).Create(claim).Error
}

func (r *repo) GetActiveClaimByTaskAndUser(ctx context.Context, taskID, userID uint64) (*CatTaskClaim, error) {
	var claim CatTaskClaim
	err := r.db.WithContext(ctx).
		Where("task_id = ? AND user_id = ? AND status = ?", taskID, userID, "claimed").
		First(&claim).Error
	if err != nil {
		return nil, err
	}
	return &claim, nil
}

func (r *repo) UpdateCatTaskClaim(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&CatTaskClaim{}).Where("id = ?", id).Updates(updates).Error
}

// ── CatAdoption ───────────────────────────────────────────────────────────────

func (r *repo) UpsertCatAdoption(ctx context.Context, catID uint64, status string) error {
	return r.db.WithContext(ctx).
		Where(CatAdoption{CatID: catID}).
		Assign(CatAdoption{Status: status}).
		FirstOrCreate(&CatAdoption{}).Error
}

// ── CatTaskProgress ───────────────────────────────────────────────────────────

func (r *repo) CreateCatTaskProgress(ctx context.Context, progress *CatTaskProgress) error {
	return r.db.WithContext(ctx).Create(progress).Error
}
