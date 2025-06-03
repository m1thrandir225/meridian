package services

import "github.com/m1thrandir225/meridian/internal/identity/infrastructure/persistence"

type UserService struct {
	repo persistence.UserRepository
}

func NewUserService(repository persistence.UserRepository) *UserService {
	return &UserService{
		repo: repository,
	}
}
