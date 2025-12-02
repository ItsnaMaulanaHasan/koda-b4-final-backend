package services

import (
	"backend-koda-shortlink/internal/models"
	"backend-koda-shortlink/internal/repository"
	"context"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetById(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetById(ctx, id)
}
