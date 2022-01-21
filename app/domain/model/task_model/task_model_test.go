package task_model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	t.Parallel()

	id := TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")

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
			time.Date(2022, 1, 23, 00, 00, 00, 000000000, time.UTC),
			&Task{taskID: id, name: "Venue Reservation", detail: "Reserve venue for conference", status: Working, completionDate: time.Time{}, deadline: time.Date(2022, 1, 23, 00, 00, 00, 000000000, time.UTC), notificationCount: 0, postponedCount: 0},
			"",
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
