package persistence

import (
	"todo-app/domain/model"
	"todo-app/domain/repository"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TaskPersistence struct {
	conn *gorm.DB
}

func NewTaskPersistence(conn *gorm.DB) repository.TaskRepository {
	return &TaskPersistence{
		conn,
	}
}

func (tp *TaskPersistence) Create(task *model.Task) error {
	if err := tp.conn.Create(&task).Error; err != nil {
		return errors.Wrapf(err, "failed to create task. task: %+v", &task)
	}

	return nil
}

func (tp *TaskPersistence) FindByID(id model.TaskID) (*model.Task, error) {
	t := &model.Task{ID: id}

	if err := tp.conn.First(&t).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find task. id: %+v", id)
	}

	return t, nil
}

func (tp *TaskPersistence) Update(t *model.Task) error {
	return tp.conn.Save(&t).Error
}
