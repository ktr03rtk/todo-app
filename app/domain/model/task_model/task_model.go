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

func CreateTask(id TaskID, name, detail string, deadline time.Time) (*Task, error) {
	t := Task{
		taskID: id,
		// taskID:            TaskID(uuid.Must(uuid.NewRandom()).String()),
		name:              name,
		detail:            detail,
		status:            Working,
		completionDate:    time.Time{},
		deadline:          deadline,
		notificationCount: 0,
		postponedCount:    0,
	}

	if !TaskSpecSatisfied(t) {
		return nil, errors.Errorf("failed to specify Task. t: %+v", &t)
	}

	return &t, nil
}

func TaskSpecSatisfied(t Task) bool {
	return (t.notificationCount <= NOTIFICATION_COUNT_LIMIT && t.postponedCount <= POSTPONED_COUNT_LIMIT)
}
