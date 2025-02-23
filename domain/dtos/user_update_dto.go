package dtos

type UserUpdateDto struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email" validate:"omitempty,email,ascii"`
	Role     string `json:"role"`
	State    string `json:"state"`
}
