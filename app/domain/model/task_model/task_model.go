package task_model

import (
	"time"

	"github.com/pkg/errors"
)

type Task struct {
	taskID            TaskID
	name              string
	detail            string
	status            Status
	completionDate    time.Time
	deadline          time.Time
	notificationCount int
	postponedCount    int
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
		taskID:            id,
		name:              name,
		detail:            detail,
		status:            Working,
		completionDate:    time.Time{},
		deadline:          deadline,
		notificationCount: 0,
		postponedCount:    0,
	}

	if err := TaskSpecSatisfied(t); err != nil {
		return nil, errors.Wrapf(err, "failed to specify Task. t: %+v", t)
	}

	return &t, nil
}

func TaskSpecSatisfied(t Task) error {
	now := getNow()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if t.deadline.Before(today) {
		return errors.Errorf("past day is set on deadline. t.deadline: %+v", t.deadline)
	}

	if t.notificationCount > NOTIFICATION_COUNT_LIMIT {
		return errors.Errorf("notification counts exceeds limit. t.notificationCount: %+v", t.notificationCount)
	}

	if t.postponedCount > POSTPONED_COUNT_LIMIT {
		return errors.Errorf("postponed counts exceeds limit. t.notificationCount: %+v", t.notificationCount)
	}

	return nil
}
