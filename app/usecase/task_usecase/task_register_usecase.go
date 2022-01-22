package task_usecase

import (
	"time"

	"github.com/pkg/errors"
	"todo-app/domain/model/task_model"
	"todo-app/domain/model/uuid_model"
	"todo-app/domain/repository/task_repository"
)

type TaskRegisterUsecase interface {
	Execute(name, detail string, deadline time.Time) error
}

type taskRegisterUsecase struct {
	taskRepository task_repository.TaskRepository
}

func NewTaskRegisterUsecase(tr task_repository.TaskRepository) TaskRegisterUsecase {
	return &taskRegisterUsecase{taskRepository: tr}
}

func (u *taskRegisterUsecase) Execute(name, detail string, deadline time.Time) error {
	id := uuid_model.CreateUUID()

	t, err := task_model.CreateTask(task_model.TaskID(id), name, detail, deadline)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	if err := u.taskRepository.Insert(*t); err != nil {
		return errors.Wrap(err, "failed to store task")
	}

	return nil
}
