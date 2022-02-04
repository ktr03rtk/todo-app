package usecase

import (
	"todo-app/domain/model"
	"todo-app/domain/repository"
	"todo-app/domain/service"

	"github.com/pkg/errors"
)

type UserUsecase interface {
	Signup(email, password string) error
	Authenticate(email, password string) error
}

type userUsecase struct {
	userRepository repository.UserRepository
	userService    service.UserService
}

func NewUserUsecase(ur repository.UserRepository, us service.UserService) UserUsecase {
	return &userUsecase{
		userRepository: ur,
		userService:    us,
	}
}

func (u *userUsecase) Signup(email, password string) error {
	ok, err := u.userService.IsExists(model.Email(email))
	if ok {
		return errors.Errorf("already registered email. email: %s", email)
	} else if err != nil {
		return err
	}

	id := model.CreateUUID()

	user, err := model.NewUser(model.UserID(id), model.Email(email), password)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	if err := u.userRepository.Create(user); err != nil {
		return errors.Wrap(err, "failed to store user")
	}

	return nil
}

func (u *userUsecase) Authenticate(email, password string) error {
	user, err := u.userRepository.FindByEmail(model.Email(email))
	if err != nil {
		return errors.Wrap(err, "failed to find user")
	}

	if err := user.ValidatePassword(password); err != nil {
		return err
	}

	return nil
}
