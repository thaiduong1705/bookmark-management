package user

import (
	"bookmark-management/internal/domain"
	"context"
	"fmt"
)

// CreateUser hashes the password and stores the user.
func (u *service) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// gen pass
	hashed, err := u.passwordUtil.Hash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user.Password = hashed
	created, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return created, nil
}
