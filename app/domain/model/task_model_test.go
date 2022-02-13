package model

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	t.Parallel()

	id := TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")
	userID := UserID("477ecd7f-48fe-6b1c-499a-ec9f52b15a33")

	getNow = func() time.Time {
		return time.Date(2022, 1, 25, 10, 10, 10, 0, time.Local)
	}

	tests := []struct {
		name           string
		taskName       string
		detail         string
		deadline       time.Time
		expectedOutput *Task
		expectedErr    error
	}{
		{
			"normal case",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Date(2022, 1, 26, 0, 0, 0, 0, time.Local),
			&Task{ID: id, UserID: userID, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: Working, CompletionDate: nil, Deadline: time.Date(2022, 1, 26, 0, 0, 0, 0, time.Local), NotificationCount: 0, PostponedCount: 0},
			nil,
		},
		{
			"normal case(on the day of deadline)",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Date(2022, 1, 25, 0, 0, 0, 0, time.Local),
			&Task{ID: id, UserID: userID, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: Working, CompletionDate: nil, Deadline: time.Date(2022, 1, 25, 0, 0, 0, 0, time.Local), NotificationCount: 0, PostponedCount: 0},
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := NewTask(id, userID, tt.taskName, tt.detail, tt.deadline)
			if err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
				assert.Exactly(t, tt.expectedOutput, output)
			}
		})
	}
}

func TestTaskSpecSatisfied(t *testing.T) {
	t.Parallel()

	id := TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")
	userID := UserID("477ecd7f-48fe-6b1c-499a-ec9f52b15a33")

	getNow = func() time.Time {
		return time.Date(2022, 1, 25, 10, 10, 10, 0, time.Local)
	}

	tests := []struct {
		name        string
		input       Task
		expectedErr error
	}{
		{
			"normal case: count is under the limit",
			Task{ID: id, UserID: userID, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: Working, CompletionDate: nil, Deadline: time.Date(2022, 1, 26, 0, 0, 0, 0, time.Local), NotificationCount: 5, PostponedCount: 3},
			nil,
		},
		{
			"error case: notification counts exceeds limit",
			Task{ID: id, UserID: userID, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: Working, CompletionDate: nil, Deadline: time.Date(2022, 1, 26, 0, 0, 0, 0, time.Local), NotificationCount: 6, PostponedCount: 0},
			errors.New("notification counts exceeds limit"),
		},
		{
			"error case: postponed counts exceeds limit",
			Task{ID: id, UserID: userID, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: Working, CompletionDate: nil, Deadline: time.Date(2022, 1, 26, 0, 0, 0, 0, time.Local), NotificationCount: 0, PostponedCount: 4},
			errors.New("postponed counts exceeds limit"),
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := TaskSpecSatisfied(tt.input); err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
			}
		})
	}
}

func TestTaskSet(t *testing.T) {
	t.Parallel()

	referenceDate := time.Date(2022, 1, 25, 10, 10, 10, 0, time.Local)
	completedDate := time.Date(2022, 1, 25, 0, 0, 0, 0, time.Local)
	behindDate := time.Date(2022, 1, 27, 10, 10, 10, 0, time.Local)
	behindCompletedDate := time.Date(2022, 1, 27, 0, 0, 0, 0, time.Local)

	createdDate := time.Date(2022, 1, 26, 0, 0, 0, 0, time.Local)
	postponedDate := time.Date(2022, 1, 27, 0, 0, 0, 0, time.Local)

	id := TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")
	userID := UserID("477ecd7f-48fe-6b1c-499a-ec9f52b15a33")
	fetchedTask := Task{ID: id, UserID: userID, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: Working, CompletionDate: nil, Deadline: createdDate, NotificationCount: 0, PostponedCount: 0}

	updatedTaskName := "Updated Venue Reservation"
	updatedTaskDetail := "Updated Reserve venue for conference"

	tests := []struct {
		name           string
		status         Status
		deadline       time.Time
		date           time.Time
		expectedOutput *Task
		expectedErr    error
	}{
		{
			"normal case",
			Working,
			createdDate,
			referenceDate,
			&Task{ID: id, UserID: userID, Name: "Updated Venue Reservation", Detail: "Updated Reserve venue for conference", Status: Working, CompletionDate: nil, Deadline: createdDate, NotificationCount: 0, PostponedCount: 0},
			nil,
		},
		{
			"postponed case",
			Working,
			postponedDate,
			referenceDate,
			&Task{ID: id, UserID: userID, Name: "Updated Venue Reservation", Detail: "Updated Reserve venue for conference", Status: Working, CompletionDate: nil, Deadline: postponedDate, NotificationCount: 0, PostponedCount: 1},
			nil,
		},
		{
			"completed in time case",
			Completed,
			createdDate,
			referenceDate,
			&Task{ID: id, UserID: userID, Name: "Updated Venue Reservation", Detail: "Updated Reserve venue for conference", Status: Completed, CompletionDate: &completedDate, Deadline: createdDate, NotificationCount: 0, PostponedCount: 0},
			nil,
		},
		{
			"behind completed case",
			Completed,
			createdDate,
			behindDate,
			&Task{ID: id, UserID: userID, Name: "Updated Venue Reservation", Detail: "Updated Reserve venue for conference", Status: Completed, CompletionDate: &behindCompletedDate, Deadline: createdDate, NotificationCount: 0, PostponedCount: 0},
			nil,
		},
		{
			"behind working case",
			Working,
			createdDate,
			behindDate,
			&Task{ID: id, UserID: userID, Name: "Updated Venue Reservation", Detail: "Updated Reserve venue for conference", Status: Behind, CompletionDate: nil, Deadline: createdDate, NotificationCount: 0, PostponedCount: 0},
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			getNow = func() time.Time { return tt.date }
			output, err := TaskSet(fetchedTask, updatedTaskName, updatedTaskDetail, tt.status, tt.deadline)
			if err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
				assert.Exactly(t, tt.expectedOutput, output)
			}
		})
	}
}
