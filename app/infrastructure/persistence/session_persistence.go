package persistence

import (
	"todo-app/domain/model"
	"todo-app/usecase"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SessionPersistence struct {
	conn *gorm.DB
}

func NewSessionPersistence(conn *gorm.DB) usecase.SessionRepository {
	return &SessionPersistence{
		conn,
	}
}

func (up *SessionPersistence) Create(s *usecase.Session) error {
	if err := up.conn.Create(&s).Error; err != nil {
		return errors.Wrapf(err, "failed to create session. session id: %+v", &s.ID)
	}

	return nil
}

func (up *SessionPersistence) FindByID(id usecase.SessionID) (*usecase.Session, error) {
	s := &usecase.Session{ID: id}

	if err := up.conn.Where(&s).First(&s).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to find session. session id: %+v", id)
	}

	return s, nil
}

func (up *SessionPersistence) FindByUserID(id model.UserID) (*usecase.Session, error) {
	s := &usecase.Session{UserID: id}

	if err := up.conn.Where(&s).First(&s).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to find session. user id: %+v", id)
	}

	return s, nil
}

func (up *SessionPersistence) Delete(id model.UserID) error {
	s := &usecase.Session{UserID: id}

	if err := up.conn.Where(&s).Delete(&s).Error; err != nil {
		return errors.Wrapf(err, "failed to create session. session id: %+v", id)
	}

	return nil
}
