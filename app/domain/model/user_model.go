package model

import (
	"regexp"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       UserID
	Email    Email
	Password string
}

type (
	UserID string
	Email  string
)

var (
	emailValidater  = regexp.MustCompile(`^.+@.+$`)
	digitValidater  = regexp.MustCompile(`\d`)
	letterValidater = regexp.MustCompile(`[a-zA-Z]`)
)

const minimumPasswordLength = 8

func NewUser(id UserID, email Email, pw string) (*User, error) {
	if err := passwordSpecSatisfied(pw); err != nil {
		return nil, errors.Wrapf(err, "invalid password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt password")
	}

	u := &User{
		ID:       id,
		Email:    email,
		Password: string(hash),
	}

	if err := UserSpecSatisfied(*u); err != nil {
		return nil, errors.Wrapf(err, "fail to satisfy User spec")
	}

	return u, nil
}

func passwordSpecSatisfied(pw string) error {
	if !digitValidater.MatchString(pw) || !letterValidater.MatchString(pw) {
		return errors.Errorf("password must contains at least one digit and letter")
	}

	if len(pw) < minimumPasswordLength {
		return errors.Errorf("password must contains at least eight characters")
	}

	return nil
}

func UserSpecSatisfied(u User) error {
	// TODO: email が未使用
	if !emailValidater.MatchString(string(u.Email)) {
		return errors.Errorf("invalid email pattern")
	}

	return nil
}
