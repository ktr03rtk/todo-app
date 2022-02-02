package persistence

import (
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserPersistence struct {
	conn *gorm.DB
}

func NewUserPersistence(conn *gorm.DB) repository.UserRepository {
	return &UserPersistence{
		conn,
	}
}

func (up *UserPersistence) Create(user *model.User) error {
	if err := up.conn.Create(&user).Error; err != nil {
		return errors.Wrapf(err, "failed to create user. user email: %+v", &user.Email)
	}

	return nil
}

func (up *UserPersistence) FindByEmail(email model.Email) (*model.User, error) {
	t := &model.User{Email: email}

	if err := up.conn.First(&t).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find user. user email: %+v", t.Email)
	}

	return t, nil
}
