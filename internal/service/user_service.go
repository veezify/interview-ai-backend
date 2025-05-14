package service

import (
	"github.com/veezify/interview-ai-backend/internal/domain/model"
	"github.com/veezify/interview-ai-backend/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.userRepository.FindByID(id)
}
