package usecase

import (
	"testing"
	"time"
	"todo-app/domain/model"
	"todo-app/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskFetchByIDUseCase(t *testing.T) {
	id := model.TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")

	tests := []struct {
		name           string
		taskID         model.TaskID
		expectedOutput *model.Task
		expectedErr    error
	}{
		{
			"normal case",
			id,
			&model.Task{ID: id, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: time.Date(2022, 1, 26, 0o0, 0o0, 0o0, 0o00000000, time.Local), NotificationCount: 0, PostponedCount: 0},
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepository := mock.NewMockTaskRepository(ctrl)
			usecase := NewTaskFetchByIDUsecase(taskRepository)

			taskRepository.EXPECT().FindByID(tt.taskID).Return(tt.expectedOutput, tt.expectedErr).Times(1)

			output, err := usecase.Execute(tt.taskID)
			if err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr)
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
