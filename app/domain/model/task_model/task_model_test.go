package task_model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	t.Parallel()

	id := TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")

	getNow = func() time.Time {
		return time.Date(2022, 1, 25, 10, 10, 10, 000000000, time.Local)
	}

	tests := []struct {
		name           string
		taskName       string
		detail         string
		deadline       time.Time
		expectedOutput *Task
		expectedErr    string
	}{
		{
			"normal case",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Date(2022, 1, 26, 00, 00, 00, 000000000, time.Local),
			&Task{taskID: id, name: "Venue Reservation", detail: "Reserve venue for conference", status: Working, completionDate: time.Time{}, deadline: time.Date(2022, 1, 26, 00, 00, 00, 000000000, time.Local), notificationCount: 0, postponedCount: 0},
			"",
		},
		{
			"normal case(on the day of deadline)",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Date(2022, 1, 25, 00, 00, 00, 000000000, time.Local),
			&Task{taskID: id, name: "Venue Reservation", detail: "Reserve venue for conference", status: Working, completionDate: time.Time{}, deadline: time.Date(2022, 1, 25, 00, 00, 00, 000000000, time.Local), notificationCount: 0, postponedCount: 0},
			"",
		},
		{
			"error case(deadline is past)",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Date(2022, 1, 24, 00, 00, 00, 000000000, time.Local),
			&Task{taskID: id, name: "Venue Reservation", detail: "Reserve venue for conference", status: Working, completionDate: time.Time{}, deadline: time.Date(2022, 1, 24, 00, 00, 00, 000000000, time.Local), notificationCount: 0, postponedCount: 0},
			"past day is set on deadline",
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := CreateTask(id, tt.taskName, tt.detail, tt.deadline)
			if err != nil {
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, "", "error is expected but received nil")
				assert.Exactly(t, tt.expectedOutput, output)
			}
		})
	}
}

func TestTaskSpecSatisfied(t *testing.T) {
	t.Parallel()

	id := TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")

	getNow = func() time.Time {
		return time.Date(2022, 1, 25, 10, 10, 10, 000000000, time.Local)
	}

	tests := []struct {
		name        string
		input       Task
		expectedErr string
	}{
		{
			"normal case: count is under the limit",
			Task{taskID: id, name: "Venue Reservation", detail: "Reserve venue for conference", status: Working, completionDate: time.Time{}, deadline: time.Date(2022, 1, 26, 00, 00, 00, 000000000, time.Local), notificationCount: 5, postponedCount: 3},
			"",
		},
		{
			"error case: notification counts exceeds limit",
			Task{taskID: id, name: "Venue Reservation", detail: "Reserve venue for conference", status: Working, completionDate: time.Time{}, deadline: time.Date(2022, 1, 26, 00, 00, 00, 000000000, time.Local), notificationCount: 6, postponedCount: 0},
			"notification counts exceeds limit",
		},
		{
			"error case: postponed counts exceeds limit",
			Task{taskID: id, name: "Venue Reservation", detail: "Reserve venue for conference", status: Working, completionDate: time.Time{}, deadline: time.Date(2022, 1, 26, 00, 00, 00, 000000000, time.Local), notificationCount: 0, postponedCount: 4},
			"postponed counts exceeds limit",
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := TaskSpecSatisfied(tt.input); err != nil {
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, "", "error is expected but received nil")
			}
		})
	}
}