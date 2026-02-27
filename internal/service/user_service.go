package service

import (
	"talk-backend/internal/models"
	"talk-backend/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetMe(userID string) (*models.User, error) {
	return s.repo.FindByID(userID)
}

func (s *UserService) RegisterUser(user *models.User) error {
	return s.repo.Create(user)
}
