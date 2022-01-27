package usecase

import (
	"fmt"
	"testing"
	"time"
	"todo-app/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreateUseCase(t *testing.T) {
	tests := []struct {
		name              string
		taskName          string
		detail            string
		deadline          time.Time
		expectedOutput    error
		expectedErr       string
		expectedCallTimes int
	}{
		{
			"normal case",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Now().AddDate(0, 0, 1),
			nil,
			"",
			1,
		},
		{
			"error case(deadline is past)",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Now().AddDate(0, 0, -1),
			nil,
			"past day is set on deadline",
			0,
		},
		{
			"error case(fail to store)",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Now().AddDate(0, 0, 1),
			fmt.Errorf("fail to store"),
			"fail to store",
			1,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepository := mock.NewMockTaskRepository(ctrl)
			usecase := NewTaskCreateUsecase(taskRepository)

			taskRepository.EXPECT().Create(gomock.Any()).Return(tt.expectedOutput).Times(tt.expectedCallTimes)

			if err := usecase.Execute(tt.taskName, tt.detail, tt.deadline); err != nil {
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.Exactly(t, tt.expectedErr, "", "error is expected but received nil")
			}
		})
	}
}
