package usecase

import (
	"time"
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
)

type TaskCreateUsecase interface {
	Execute(name, detail string, deadline time.Time) error
}

type taskCreateUsecase struct {
	taskRepository repository.TaskRepository
}

func NewTaskCreateUsecase(tr repository.TaskRepository) TaskCreateUsecase {
	return &taskCreateUsecase{taskRepository: tr}
}

func (u *taskCreateUsecase) Execute(name, detail string, deadline time.Time) error {
	id := model.CreateUUID()

	t, err := model.NewTask(model.TaskID(id), name, detail, deadline)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	if err := u.taskRepository.Create(t); err != nil {
		return errors.Wrap(err, "failed to store task")
	}

	return nil
}
