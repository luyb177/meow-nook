package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository defines the data access interface for users.
type Repository interface {
	// Create inserts a new user and returns the created record (with ID filled).
	Create(ctx context.Context, u *User) error

	// ExistsByEmail returns true if a user with the given email already exists.
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// FindByEmail returns the user with the given email, or nil if not found.
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByID returns the user with the given id, or nil if not found.
	FindByID(ctx context.Context, userID int64) (*User, error)

	// UpdateFields updates the given columns for the target user id.
	UpdateFields(ctx context.Context, userID int64, fields map[string]any) error

	// AddPointsDelta adds delta to points atomically and returns new points.
	AddPointsDelta(ctx context.Context, userID int64, delta int32) (int32, error)
}

type repo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *repo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repo) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) FindByID(ctx context.Context, userID int64) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) UpdateFields(ctx context.Context, userID int64, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(fields).Error
}

func (r *repo) AddPointsDelta(ctx context.Context, userID int64, delta int32) (int32, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer func() { _ = tx.Rollback() }()

	var u User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&u).Error; err != nil {
		return 0, err
	}

	newPoints := u.Points + delta
	if newPoints < 0 {
		newPoints = 0
	}

	if err := tx.Model(&User{}).Where("id = ?", userID).Update("points", newPoints).Error; err != nil {
		return 0, err
	}
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return newPoints, nil
}
