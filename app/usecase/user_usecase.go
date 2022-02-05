package usecase

import (
	"todo-app/domain/model"
	"todo-app/domain/repository"
	"todo-app/domain/service"

	"github.com/pkg/errors"
)

type UserUsecase interface {
	SignUp(email, password string) error
	Authenticate(email, password string) (model.UserID, error)
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

func (u *userUsecase) SignUp(email, password string) error {
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

func (u *userUsecase) Authenticate(email, password string) (model.UserID, error) {
	user, err := u.userRepository.FindByEmail(model.Email(email))
	if err != nil {
		return "", errors.Wrap(err, "failed to find user")
	} else if user == nil {
		return "", errors.New("user is not registered")
	}

	if err := user.ValidatePassword(password); err != nil {
		return "", err
	}

	return user.ID, nil
}
