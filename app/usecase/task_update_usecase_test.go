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

func TestTaskUpdateUseCase(t *testing.T) {
	id := model.TaskID("19742914-f296-4855-aa8d-f099727e288f")
	dl := time.Now().AddDate(0, 0, 2)
	deadline := time.Date(dl.Year(), dl.Month(), dl.Day(), 0, 0, 0, 0, time.Local)

	normalTask := &model.Task{ID: id, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: deadline, NotificationCount: 0, PostponedCount: 0}
	updatedTask := &model.Task{ID: id, Name: "Updated Venue Reservation", Detail: "Updated Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: deadline, NotificationCount: 0, PostponedCount: 0}

	updatedTaskName := "Updated Venue Reservation"
	updatedTaskDetail := "Updated Reserve venue for conference"
	status := model.Working

	tests := []struct {
		name                string
		expectedFindByIDErr error
		expectedUpdateErr   error
		expectedErr         error
		expectedCallTimes   int
	}{
		{
			"normal case",
			nil,
			nil,
			nil,
			1,
		},
		{
			"find by id error case",
			errors.New("find by id error"),
			nil,
			errors.New("failed to fetch task"),
			0,
		},
		{
			"update error case",
			nil,
			errors.New("update error"),
			errors.New("failed to update task"),
			1,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepository := mock.NewMockTaskRepository(ctrl)
			usecase := NewTaskUpdateUsecase(taskRepository)

			gomock.InOrder(
				taskRepository.EXPECT().FindByID(id).Return(normalTask, tt.expectedFindByIDErr).Times(1),
				taskRepository.EXPECT().Update(updatedTask).Return(tt.expectedUpdateErr).Times(tt.expectedCallTimes),
			)

			if err := usecase.Execute(id, updatedTaskName, updatedTaskDetail, status, deadline); err != nil {
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
			}
		})
	}
}

func TestTaskUpdatePostoneUseCase(t *testing.T) {
	id := model.TaskID("19742914-f296-4855-aa8d-f099727e288f")
	dl := time.Now().AddDate(0, 0, 2)
	deadline := time.Date(dl.Year(), dl.Month(), dl.Day(), 0, 0, 0, 0, time.Local)
	updatedDeadline := deadline.Add(24 * time.Hour)

	updatedTask := &model.Task{ID: id, Name: "Updated Venue Reservation", Detail: "Updated Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: updatedDeadline, NotificationCount: 0, PostponedCount: 1}

	updatedTaskName := "Updated Venue Reservation"
	updatedTaskDetail := "Updated Reserve venue for conference"
	status := model.Working

	tests := []struct {
		name              string
		fetchedTask       *model.Task
		expectedErr       error
		expectedCallTimes int
	}{
		{
			"normal case",
			&model.Task{ID: id, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: deadline, NotificationCount: 0, PostponedCount: 0},
			nil,
			1,
		},
		{
			"postponed count limit over erro case",
			&model.Task{ID: id, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: deadline, NotificationCount: 0, PostponedCount: 3},
			errors.New("failed to satisfy task spec"),
			0,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepository := mock.NewMockTaskRepository(ctrl)
			usecase := NewTaskUpdateUsecase(taskRepository)

			gomock.InOrder(
				taskRepository.EXPECT().FindByID(id).Return(tt.fetchedTask, nil).Times(1),
				taskRepository.EXPECT().Update(updatedTask).Return(nil).Times(tt.expectedCallTimes),
			)

			if err := usecase.Execute(id, updatedTaskName, updatedTaskDetail, status, updatedDeadline); err != nil {
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
			}
		})
	}
}
