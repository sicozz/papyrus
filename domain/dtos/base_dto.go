package dtos

import (
	"fmt"
)

type BaseDto struct {
	Message string `json:"message"`
}

func NewBaseDto(msgs ...string) (dto BaseDto) {
	dto = BaseDto{fmt.Sprintln(msgs)}
	return
}
