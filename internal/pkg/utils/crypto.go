package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashStr(source string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(source), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func UuidStr() uuid.UUID {
	uuid := uuid.New()
	return uuid
}
