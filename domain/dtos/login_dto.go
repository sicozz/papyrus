package dtos

type LoginDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewLoginDto(username string, password string) (dto LoginDto) {
	dto = LoginDto{username, password}
	return
}
