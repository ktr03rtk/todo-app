//go:generate mockgen -source=task_repository.go -destination=../../mock/mock_task_repository.go -package=mock
package repository

import "todo-app/domain/model"

type UserRepository interface {
	Create(*model.User) error
	FindByID(model.Email) (*model.User, error)
}
