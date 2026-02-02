package service

import (
	"talk-backend/internal/models"
	"talk-backend/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func (s *UserService) RegisterUser(user *models.User) error {
	return s.repo.Create(user)
}
