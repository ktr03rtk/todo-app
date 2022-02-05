package usecase

import (
	"errors"
	"testing"
	"todo-app/domain/model"
	"todo-app/domain/service"
	"todo-app/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserSignupUseCase(t *testing.T) {
	id := model.UserID("72c24944-f532-4c5d-a695-70fa3e72f3ab")

	tests := []struct {
		name                 string
		email                string
		password             string
		findByEmailOutput    *model.User
		findByEmailErrOutput error
		CreateErrOutput      error
		expectedOutput       error
		expectedCallTimes    int
	}{
		{
			"normal case",
			"abc@example.com",
			"password123",
			nil,
			nil,
			nil,
			nil,
			1,
		},
		{
			"failed to find user case",
			"abc@example.com",
			"password123",
			nil,
			errors.New("fail find by email"),
			nil,
			errors.New("failed to find user"),
			0,
		},
		{
			"already user registered case",
			"abc@example.com",
			"password123",
			&model.User{ID: id, Email: "abc@example.com", Password: "$2a$10$bUJO2D0iREJl.350fkaJIeXVdEL9yNcHT8smkC90j0kQ9okVVKfsq"},
			nil,
			nil,
			errors.New("already registered email"),
			0,
		},
		{
			"failed to create user case",
			"abc@example.com",
			"password123",
			nil,
			nil,
			errors.New("fail create"),
			errors.New("failed to store user"),
			1,
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepository := mock.NewMockUserRepository(ctrl)
			userService := service.NewUService(userRepository)
			usecase := NewUserUsecase(userRepository, userService)

			gomock.InOrder(
				userRepository.EXPECT().FindByEmail(model.Email(tt.email)).Return(tt.findByEmailOutput, tt.findByEmailErrOutput).Times(1),
				userRepository.EXPECT().Create(gomock.Any()).Return(tt.CreateErrOutput).Times(tt.expectedCallTimes),
			)

			if err := usecase.SignUp(tt.email, tt.password); err != nil {
				assert.Contains(t, err.Error(), tt.expectedOutput.Error())
			} else {
				assert.Nil(t, tt.expectedOutput, "error is expected but received nil")
			}
		})
	}
}
