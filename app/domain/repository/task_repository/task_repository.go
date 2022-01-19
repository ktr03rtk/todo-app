//go:generate mockgen -source=task_repository.go -destination=../../../mock/mock_task_repository.go -package=mock
package task_repository

import (
	"todo-app/domain/model/task_model"
)

type TaskRepository interface {
	Insert(task_model.Task) error
	FindByID(task_model.TaskID) (task_model.Task, error)
}
