package usecase

import (
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
)

type TaskFetchByIDUsecase interface {
	Execute(id model.TaskID) (*model.Task, error)
}

type taskFetchByIDUsecase struct {
	taskRepository repository.TaskRepository
}

func NewTaskFetchByIDUsecase(tr repository.TaskRepository) TaskFetchByIDUsecase {
	return &taskFetchByIDUsecase{taskRepository: tr}
}

func (u *taskFetchByIDUsecase) Execute(id model.TaskID) (*model.Task, error) {
	t, err := u.taskRepository.FindByID(id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch task, taskID: %s", id)
	}

	return t, nil
}
