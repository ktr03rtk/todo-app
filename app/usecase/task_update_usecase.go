package usecase

import (
	"time"
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
)

type TaskUpdateUsecase interface {
	Execute(id model.TaskID, name, detail string, status model.Status, deadline time.Time) error
}

type taskUpdateUsecase struct {
	taskRepository repository.TaskRepository
}

func NewTaskUpdateUsecase(tr repository.TaskRepository) TaskUpdateUsecase {
	return &taskUpdateUsecase{taskRepository: tr}
}

func (u *taskUpdateUsecase) Execute(id model.TaskID, name, detail string, status model.Status, deadline time.Time) error {
	fetchedTask, err := u.taskRepository.FindByID(id)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch task, taskID: %s", id)
	}

	t, err := model.TaskSet(*fetchedTask, name, detail, status, deadline)
	if err != nil {
		return errors.Wrap(err, "failed to set task")
	}

	if err := model.TaskSpecSatisfied(*t); err != nil {
		return errors.Wrap(err, "failed to satisfy task spec")
	}

	if err := u.taskRepository.Update(t); err != nil {
		return errors.Wrap(err, "failed to update task")
	}

	return nil
}
