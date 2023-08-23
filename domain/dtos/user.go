package dtos

type UserGetDto struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username" validate:"required,ascii"`
	Email    string `json:"email" validate:"required,email,ascii"`
	Name     string `json:"name" validate:"required,ascii"`
	Lastname string `json:"lastname" validate:"required,ascii"`
	// Role     string `json:"role"`
	Role  RoleDto   `json:"role"`
	State UStateDto `json:"state"`
}

type RoleDto struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
}

type UStateDto struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
}

type UserUpdateDto struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email" validate:"omitempty,email,ascii"`
	Role     string `json:"role"`
	State    string `json:"state"`
}

type UserChgPasswdDto struct {
	Passwd    string `json:"password" validate:"required"`
	NPasswd   string `json:"new_password" validate:"required"`
	ReNPasswd string `json:"re_new_password" validate:"required"`
}
