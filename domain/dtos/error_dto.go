package dtos

type ErrorDto struct {
	Message string `json:"error"`
}

func NewErrorDto(msg string) (dto ErrorDto) {
	dto = ErrorDto{msg}
	return
}
