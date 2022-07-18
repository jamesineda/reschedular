package utils

import uuid "github.com/satori/go.uuid"

type IdGenny interface {
	GenerateId() string
}

type UUIDID struct {
}

func GenerateId() string {
	return uuid.NewV4().String()
}
