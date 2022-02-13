//go:generate mockgen -source=user_repository.go -destination=../../mock/mock_user_repository.go -package=mock
package repository

import "todo-app/domain/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(model.Email) (*model.User, error)
}
