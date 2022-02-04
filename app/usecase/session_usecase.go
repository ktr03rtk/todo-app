package usecase

import (
	"time"
	"todo-app/domain/model"

	"github.com/pkg/errors"
)

type SessionUsecase interface {
	CreateSession(model.UserID) (*Session, error)
	Verify(SessionID) (*Session, error)
	DeleteSession(SessionID) error
}

type sessionUsecase struct {
	sessionRepository SessionRepository
}

type SessionRepository interface {
	Create(*Session) error
	FindByID(SessionID) (*Session, error)
	FindByUserID(model.UserID) (*Session, error)
	Delete(SessionID) error
}

func NewSessionUsecase(r SessionRepository) SessionUsecase {
	return &sessionUsecase{
		sessionRepository: r,
	}
}

type Session struct {
	ID        SessionID
	UserID    model.UserID
	CreatedAt time.Time
	ExpiredAt time.Time
}

type SessionID string

const sessionValidDuration = 2 * time.Hour

var getNow = time.Now

func (u *sessionUsecase) CreateSession(userID model.UserID) (*Session, error) {
	id := model.CreateUUID()

	now := getNow()

	s := &Session{
		ID:        SessionID(id),
		UserID:    userID,
		CreatedAt: now,
		ExpiredAt: now.Add(sessionValidDuration),
	}

	if err := u.sessionRepository.Create(s); err != nil {
		return nil, errors.Wrap(err, "failed to store session")
	}

	return s, nil
}

func (u *sessionUsecase) Verify(id SessionID) (*Session, error) {
	s, err := u.sessionRepository.FindByID(id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find session, sessionID: %s", id)
	} else if s == nil {
		return nil, nil
	}

	if getNow().After(s.ExpiredAt) {
		if err := u.DeleteSession(s.ID); err != nil {
			return nil, errors.Wrapf(err, "failed to delete expired session, sessionID: %s", s.ID)
		}

		return nil, nil
	}

	return s, nil
}

func (u *sessionUsecase) DeleteSession(id SessionID) error {
	if err := u.sessionRepository.Delete(id); err != nil {
		return errors.Wrapf(err, "failed to delete session, sessionID: %s", id)
	}

	return nil
}
