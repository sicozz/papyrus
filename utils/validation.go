package utils

import (
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsRequestValid(p any) (bool, error) {
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		return false, err
	}
	return true, nil
}
