package taskUsecase

import (
	"time"
	"todo-app/domain/model/taskModel"
	"todo-app/domain/model/uuidModel"
	"todo-app/domain/repository/taskRepository"

	"github.com/pkg/errors"
)

type TaskRegisterUsecase interface {
	Execute(name, detail string, deadline time.Time) error
}

type taskRegisterUsecase struct {
	taskRepository taskRepository.TaskRepository
}

func NewTaskRegisterUsecase(tr taskRepository.TaskRepository) TaskRegisterUsecase {
	return &taskRegisterUsecase{taskRepository: tr}
}

func (u *taskRegisterUsecase) Execute(name, detail string, deadline time.Time) error {
	id := uuidModel.CreateUUID()

	t, err := taskModel.CreateTask(taskModel.TaskID(id), name, detail, deadline)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	if err := u.taskRepository.Create(*t); err != nil {
		return errors.Wrap(err, "failed to store task")
	}

	return nil
}
