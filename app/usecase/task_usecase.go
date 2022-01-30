package usecase

import (
	"time"
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
)

type TaskUsecase interface {
	Create(name, detail string, deadline time.Time) error
	FindByID(id model.TaskID) (*model.Task, error)
	FindAll() ([]*model.Task, error)
	Update(id model.TaskID, name, detail string, status model.Status, deadline time.Time) error
}

type taskUsecase struct {
	taskRepository repository.TaskRepository
}

func NewTaskUsecase(tr repository.TaskRepository) TaskUsecase {
	return &taskUsecase{taskRepository: tr}
}

func (u *taskUsecase) Create(name, detail string, deadline time.Time) error {
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

func (u *taskUsecase) FindByID(id model.TaskID) (*model.Task, error) {
	t, err := u.taskRepository.FindByID(id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find task, taskID: %s", id)
	}

	return t, nil
}

func (u *taskUsecase) FindAll() ([]*model.Task, error) {
	tasks, err := u.taskRepository.FindAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find all tasks")
	}

	return tasks, nil
}

func (u *taskUsecase) Update(id model.TaskID, name, detail string, status model.Status, deadline time.Time) error {
	fetchedTask, err := u.taskRepository.FindByID(id)
	if err != nil {
		return errors.Wrapf(err, "failed to find task")
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
