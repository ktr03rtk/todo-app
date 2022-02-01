package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	id := UserID("72c24944-f532-4c5d-a695-70fa3e72f3ab")

	tests := []struct {
		name           string
		email          Email
		password       string
		expectedOutput *User
		expectedErr    error
	}{
		{
			"normal case",
			"abc@example.com",
			"password123",
			&User{ID: id, Email: "abc@example.com", Password: "$2a$10$bUJO2D0iREJl.350fkaJIeXVdEL9yNcHT8smkC90j0kQ9okVVKfsq"},
			nil,
		},
		{
			"invalid password (not contains digit) case",
			"abc@example.com",
			"passwordabc",
			nil,
			errors.New("password must contains at least one digit and letter"),
		},
		{
			"invalid password (not contains letter) case",
			"abc@example.com",
			"12345678999",
			nil,
			errors.New("password must contains at least one digit and letter"),
		},
		{
			"invalid password (short length) case",
			"abc@example.com",
			"pass123",
			nil,
			errors.New("password must contains at least eight characters"),
		},
		{
			"invalid email pattertn case",
			"abcexample.com",
			"password123",
			nil,
			errors.New("invalid email pattern"),
		},
	}

	for _, tt := range tests {
		tt := tt // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := NewUser(id, tt.email, tt.password)
			if err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
				assert.Exactly(t, tt.expectedOutput.ID, output.ID)
				assert.Exactly(t, tt.expectedOutput.Email, output.Email)
				assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(output.Password), []byte(tt.password)))
				assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(tt.expectedOutput.Password), []byte(tt.password)))
			}
		})
	}
}
