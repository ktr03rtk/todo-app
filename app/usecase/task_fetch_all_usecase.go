package usecase

import (
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
)

type TaskFetchAllUsecase interface {
	Execute() ([]*model.Task, error)
}

type taskFetchAllUsecase struct {
	taskRepository repository.TaskRepository
}

func NewTaskFetchAllUsecase(tr repository.TaskRepository) TaskFetchAllUsecase {
	return &taskFetchAllUsecase{taskRepository: tr}
}

func (u *taskFetchAllUsecase) Execute() ([]*model.Task, error) {
	tasks, err := u.taskRepository.FindAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch tasks")
	}

	return tasks, nil
}
