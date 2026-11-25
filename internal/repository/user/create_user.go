package user

import (
	"bookmark-management/internal/domain"
	"bookmark-management/internal/repository/dberr"
	"context"
)

// CreateUser inserts the user into the database, translating DB errors to domain errors
// (domain.ErrAlreadyExists, domain.ErrInvalidInput).
func (r *repository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, dberr.Translate(err)
	}
	return user, nil
}
