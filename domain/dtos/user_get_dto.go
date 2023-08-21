package dtos

type UserGetDto struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username" validate:"required,ascii"`
	Email    string `json:"email" validate:"required,email,ascii"`
	Name     string `json:"name" validate:"required,ascii"`
	Lastname string `json:"lastname" validate:"required,ascii"`
	Role     string `json:"role"`
	State    string	`json:"state"`
}
