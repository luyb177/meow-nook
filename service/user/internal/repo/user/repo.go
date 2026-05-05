package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository defines the data access interface for users.
type Repository interface {
	// Create inserts a new user and returns the created record (with ID filled).
	Create(ctx context.Context, u *User) error

	// ExistsByEmail returns true if a user with the given email already exists.
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// FindByEmail returns the user with the given email, or nil if not found.
	FindByEmail(ctx context.Context, email string) (*User, error)
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
