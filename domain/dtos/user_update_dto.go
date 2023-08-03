package dtos

type UserUpdateDto struct {
	Username string `json:"username"`
	Lastname string `json:"lastname"`
	Role     string `json:"role"`
	State    string `json:"state"`
}
