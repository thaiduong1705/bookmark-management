package user

import (
	"bookmark-management/internal/domain"
	"context"

	"gorm.io/gorm"
)

// Repository stores and loads users.
type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a GORM-backed user Repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
