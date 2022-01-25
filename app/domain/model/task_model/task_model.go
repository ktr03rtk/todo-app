package task_model

import (
	"time"

	"github.com/pkg/errors"
)

type Task struct {
	ID                TaskID
	Name              string
	Detail            string
	Status            Status
	CompletionDate    *time.Time
	Deadline          time.Time
	NotificationCount int
	PostponedCount    int
}

type TaskID string

type Status int

const (
	Working Status = iota
	Completed
	Behind
)

const (
	NOTIFICATION_COUNT_LIMIT = 5
	POSTPONED_COUNT_LIMIT    = 3
)

var getNow = time.Now

func CreateTask(id TaskID, name, detail string, deadline time.Time) (*Task, error) {
	t := Task{
		ID:                id,
		Name:              name,
		Detail:            detail,
		Status:            Working,
		CompletionDate:    nil,
		Deadline:          deadline,
		NotificationCount: 0,
		PostponedCount:    0,
	}

	if err := TaskSpecSatisfied(t); err != nil {
		return nil, errors.Wrapf(err, "failed to specify Task. t: %+v", t)
	}

	return &t, nil
}

func TaskSpecSatisfied(t Task) error {
	now := getNow()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if t.Deadline.Before(today) {
		return errors.Errorf("past day is set on deadline. t.deadline: %+v", t.Deadline)
	}

	if t.NotificationCount > NOTIFICATION_COUNT_LIMIT {
		return errors.Errorf("notification counts exceeds limit. t.notificationCount: %+v", t.NotificationCount)
	}

	if t.PostponedCount > POSTPONED_COUNT_LIMIT {
		return errors.Errorf("postponed counts exceeds limit. t.notificationCount: %+v", t.PostponedCount)
	}

	return nil
}
