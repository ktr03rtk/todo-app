package model

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

func NewTask(id TaskID, name, detail string, deadline time.Time) (*Task, error) {
	dl := time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 0, 0, 0, 0, time.Local)

	t := &Task{
		ID:                id,
		Name:              name,
		Detail:            detail,
		Status:            Working,
		CompletionDate:    nil,
		Deadline:          dl,
		NotificationCount: 0,
		PostponedCount:    0,
	}

	if err := TaskSpecSatisfied(*t); err != nil {
		return nil, errors.Wrapf(err, "failed to satisfy Task spec. t: %+v", t)
	}

	return t, nil
}

func TaskSpecSatisfied(t Task) error {
	if t.NotificationCount > NOTIFICATION_COUNT_LIMIT {
		return errors.Errorf("notification counts exceeds limit. t.notificationCount: %+v", t.NotificationCount)
	}

	if t.PostponedCount > POSTPONED_COUNT_LIMIT {
		return errors.Errorf("postponed counts exceeds limit. t.notificationCount: %+v", t.PostponedCount)
	}

	return nil
}

func TaskSet(fetchedTask Task, name, detail string, status Status, deadline time.Time) (*Task, error) {
	t, err := NewTask(fetchedTask.ID, name, detail, deadline)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set task")
	}

	t.Status = status
	t.PostponedCount = fetchedTask.PostponedCount
	t.NotificationCount = fetchedTask.NotificationCount

	if t.Deadline.After(fetchedTask.Deadline) {
		t.PostponedCount++
	}

	return calculate(*t), nil
}

func calculate(t Task) *Task {
	now := getNow()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if t.Status != Completed && today.After(t.Deadline) {
		t.Status = Behind
	}

	if t.Status == Completed && t.CompletionDate == nil {
		t.CompletionDate = &today
	}

	return &t
}
