package usecase

import (
	"errors"
	"testing"
	"time"
	"todo-app/domain/model"
	"todo-app/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskFetchAllUseCase(t *testing.T) {
	id1 := model.TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")
	id2 := model.TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ac")

	tests := []struct {
		name               string
		expectedOutput     []*model.Task
		expectedFindAllErr error
		expectedErr        error
	}{
		{
			"normal case",
			[]*model.Task{
				{ID: id1, Name: "Venue Reservation1", Detail: "Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: time.Date(2022, 1, 26, 0o0, 0o0, 0o0, 0o00000000, time.Local), NotificationCount: 0, PostponedCount: 0},
				{ID: id2, Name: "Venue Reservation2", Detail: "Reserve venue for conference2", Status: model.Working, CompletionDate: nil, Deadline: time.Date(2022, 1, 26, 0o0, 0o0, 0o0, 0o00000000, time.Local), NotificationCount: 1, PostponedCount: 1},
			},
			nil,
			nil,
		},
		{
			"error case",
			nil,
			errors.New("find all error"),
			errors.New("failed to fetch tasks"),
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepository := mock.NewMockTaskRepository(ctrl)
			usecase := NewTaskFetchAllUsecase(taskRepository)

			taskRepository.EXPECT().FindAll().Return(tt.expectedOutput, tt.expectedFindAllErr).Times(1)

			output, err := usecase.Execute()
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
