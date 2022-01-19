package task_model

import (
	"time"
)

type Task struct {
	taskID            TaskID
	name              string
	detail            string
	createdAt         time.Time
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
	Expired
)
