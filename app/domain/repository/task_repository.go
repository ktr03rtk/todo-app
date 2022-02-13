//go:generate mockgen -source=task_repository.go -destination=../../mock/mock_task_repository.go -package=mock
package repository

import "todo-app/domain/model"

type TaskRepository interface {
	Create(*model.Task) error
	FindByID(model.TaskID) (*model.Task, error)
	FindAll() ([]*model.Task, error)
	Update(*model.Task) error
}
