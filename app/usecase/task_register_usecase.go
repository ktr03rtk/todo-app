package usecase

import (
	"time"
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
)

type TaskRegisterUsecase interface {
	Execute(name, detail string, deadline time.Time) error
}

type taskRegisterUsecase struct {
	taskRepository repository.TaskRepository
}

func NewTaskRegisterUsecase(tr repository.TaskRepository) TaskRegisterUsecase {
	return &taskRegisterUsecase{taskRepository: tr}
}

func (u *taskRegisterUsecase) Execute(name, detail string, deadline time.Time) error {
	id := model.CreateUUID()

	t, err := model.CreateTask(model.TaskID(id), name, detail, deadline)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	if err := u.taskRepository.Create(t); err != nil {
		return errors.Wrap(err, "failed to store task")
	}

	return nil
}
