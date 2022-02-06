package usecase

import (
	"testing"
	"time"
	"todo-app/domain/model"
	"todo-app/mock"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreateUseCase(t *testing.T) {
	session := Session{UserID: model.UserID("477ecd7f-48fe-6b1c-499a-ec9f52b15a33")}

	tests := []struct {
		name              string
		taskName          string
		detail            string
		deadline          time.Time
		expectedOutput    error
		expectedErr       error
		expectedCallTimes int
	}{
		{
			"normal case",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Now().AddDate(0, 0, 2),
			nil,
			nil,
			1,
		},
		{
			"error case(fail to store)",
			"Venue Reservation",
			"Reserve venue for conference",
			time.Now().AddDate(0, 0, 2),
			errors.New("failed to create"),
			errors.New("failed to store"),
			1,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepository := mock.NewMockTaskRepository(ctrl)
			usecase := NewTaskUsecase(taskRepository)

			taskRepository.EXPECT().Create(gomock.Any()).Return(tt.expectedOutput).Times(tt.expectedCallTimes)

			if err := usecase.Create(session, tt.taskName, tt.detail, tt.deadline); err != nil {
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

func TestTaskFindByIDUseCase(t *testing.T) {
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
			usecase := NewTaskUsecase(taskRepository)

			taskRepository.EXPECT().FindByID(tt.taskID).Return(tt.expectedOutput, tt.expectedErr).Times(1)

			output, err := usecase.FindByID(tt.taskID)
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

func TestTaskFindAllUseCase(t *testing.T) {
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
			errors.New("failed to find all tasks"),
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepository := mock.NewMockTaskRepository(ctrl)
			usecase := NewTaskUsecase(taskRepository)

			taskRepository.EXPECT().FindAll().Return(tt.expectedOutput, tt.expectedFindAllErr).Times(1)

			output, err := usecase.FindAll()
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

func TestTaskUpdateUseCase(t *testing.T) {
	session := Session{UserID: model.UserID("477ecd7f-48fe-6b1c-499a-ec9f52b15a33")}
	id := model.TaskID("19742914-f296-4855-aa8d-f099727e288f")
	dl := time.Now().AddDate(0, 0, 2)
	deadline := time.Date(dl.Year(), dl.Month(), dl.Day(), 0, 0, 0, 0, time.Local)

	normalTask := &model.Task{ID: id, UserID: session.UserID, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: deadline, NotificationCount: 0, PostponedCount: 0}
	updatedTask := &model.Task{ID: id, UserID: session.UserID, Name: "Updated Venue Reservation", Detail: "Updated Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: deadline, NotificationCount: 0, PostponedCount: 0}

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
			errors.New("failed to find task"),
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
			usecase := NewTaskUsecase(taskRepository)

			gomock.InOrder(
				taskRepository.EXPECT().FindByID(id).Return(normalTask, tt.expectedFindByIDErr).Times(1),
				taskRepository.EXPECT().Update(updatedTask).Return(tt.expectedUpdateErr).Times(tt.expectedCallTimes),
			)

			if err := usecase.Update(session, id, updatedTaskName, updatedTaskDetail, status, deadline); err != nil {
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

func TestTaskUpdatePostoneUseCase(t *testing.T) {
	session := Session{UserID: model.UserID("477ecd7f-48fe-6b1c-499a-ec9f52b15a33")}
	id := model.TaskID("19742914-f296-4855-aa8d-f099727e288f")
	dl := time.Now().AddDate(0, 0, 2)
	deadline := time.Date(dl.Year(), dl.Month(), dl.Day(), 0, 0, 0, 0, time.Local)
	updatedDeadline := deadline.Add(24 * time.Hour)

	updatedTask := &model.Task{ID: id, UserID: session.UserID, Name: "Updated Venue Reservation", Detail: "Updated Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: updatedDeadline, NotificationCount: 0, PostponedCount: 1}

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
			&model.Task{ID: id, UserID: session.UserID, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: deadline, NotificationCount: 0, PostponedCount: 0},
			nil,
			1,
		},
		{
			"postponed count limit over erro case",
			&model.Task{ID: id, UserID: session.UserID, Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: deadline, NotificationCount: 0, PostponedCount: 3},
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
			usecase := NewTaskUsecase(taskRepository)

			gomock.InOrder(
				taskRepository.EXPECT().FindByID(id).Return(tt.fetchedTask, nil).Times(1),
				taskRepository.EXPECT().Update(updatedTask).Return(nil).Times(tt.expectedCallTimes),
			)

			if err := usecase.Update(session, id, updatedTaskName, updatedTaskDetail, status, updatedDeadline); err != nil {
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

func TestOtherUsersTaskUpdateUseCase(t *testing.T) {
	session := Session{UserID: model.UserID("477ecd7f-48fe-6b1c-499a-ec9f52b15a33")}
	id := model.TaskID("19742914-f296-4855-aa8d-f099727e288f")
	dl := time.Now().AddDate(0, 0, 2)
	deadline := time.Date(dl.Year(), dl.Month(), dl.Day(), 0, 0, 0, 0, time.Local)

	otherUsersTask := &model.Task{ID: id, UserID: model.UserID("xxxecd7f-48fe-6b1c-499a-ec9f52b15a33"), Name: "Venue Reservation", Detail: "Reserve venue for conference", Status: model.Working, CompletionDate: nil, Deadline: deadline, NotificationCount: 0, PostponedCount: 0}

	updatedTaskName := "Updated Venue Reservation"
	updatedTaskDetail := "Updated Reserve venue for conference"
	status := model.Working
	expectedErr := errors.New("session user is not task owner")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepository := mock.NewMockTaskRepository(ctrl)
	usecase := NewTaskUsecase(taskRepository)

	taskRepository.EXPECT().FindByID(id).Return(otherUsersTask, nil).Times(1)

	if err := usecase.Update(session, id, updatedTaskName, updatedTaskDetail, status, deadline); err != nil {
		if expectedErr != nil {
			assert.Contains(t, err.Error(), expectedErr.Error())
		} else {
			t.Fatalf("error is not expected but received: %v", err)
		}
	} else {
		assert.Exactly(t, expectedErr, nil, "error is expected but received nil")
	}
}
