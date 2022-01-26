package model

import "github.com/google/uuid"

type UUID string

func CreateUUID() UUID {
	return UUID(uuid.Must(uuid.NewRandom()).String())
}
