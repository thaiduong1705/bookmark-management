package user

import (
	"bookmark-management/internal/domain"
	userRepository "bookmark-management/internal/repository/user"
	"bookmark-management/pkg/utils"
	"context"
)

// Service contains the user business logic.
type Service interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type service struct {
	userRepo     userRepository.Repository
	passwordUtil utils.PasswordUtil
}

// NewService creates a user Service.
func NewService(userRepo userRepository.Repository, passwordUtil utils.PasswordUtil) Service {
	return &service{
		userRepo:     userRepo,
		passwordUtil: passwordUtil,
	}
}
