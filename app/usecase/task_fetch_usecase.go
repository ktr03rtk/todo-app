package usecase

import (
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
)

type TaskFetchUsecase interface {
	Execute(id model.TaskID) (*model.Task, error)
}

type taskFetchUsecase struct {
	taskRepository repository.TaskRepository
}

func NewTaskFetchUsecase(tr repository.TaskRepository) TaskFetchUsecase {
	return &taskFetchUsecase{taskRepository: tr}
}

func (u *taskFetchUsecase) Execute(id model.TaskID) (*model.Task, error) {
	t, err := u.taskRepository.FindByID(id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch task, taskID: %s", id)
	}

	return t, nil
}
