package dtos

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewLoginDto(username string, password string) (dto LoginDto) {
	dto = LoginDto{username, password}
	return
}
