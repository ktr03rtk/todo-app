//go:generate mockgen -source=task_repository.go -destination=../../../mock/mock_task_repository.go -package=mock
package taskRepository

import "todo-app/domain/model/taskModel"

type TaskRepository interface {
	Create(taskModel.Task) error
	FindByID(taskModel.TaskID) (taskModel.Task, error)
}
