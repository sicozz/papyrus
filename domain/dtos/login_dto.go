package dtos

type LoginDto struct {
	Username string `json:"username" validate:"required,excludesall=!@#?"`
	Password string `json:"password" validate:"required,excludesall=!@#?"`
}
