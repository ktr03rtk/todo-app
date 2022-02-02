package service

import (
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
)

type UserService interface {
	IsExists(model.Email) (bool, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUService(ur repository.UserRepository) UserService {
	return &userService{userRepository: ur}
}

func (s *userService) IsExists(email model.Email) (bool, error) {
	u, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return false, errors.Wrapf(err, "failed to find user, email: %s", email)
	}

	if u.ID == "" {
		return false, nil
	}

	return true, nil
}
