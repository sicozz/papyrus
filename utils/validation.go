package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
