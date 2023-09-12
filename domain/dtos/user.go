package dtos

type UserGetDto struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username" validate:"required,excludesall=!@#?"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,excludesall=!@#?"`
	Lastname string `json:"lastname" validate:"required,excludesall=!@#?"`
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

type UserStore struct {
	Username string `json:"username" validate:"required,excludesall=!@#?"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,excludesall=!@#?"`
	Name     string `json:"name" validate:"required,excludesall=!@#?"`
	Lastname string `json:"lastname" validate:"required,excludesall=!@#?"`
	Role     string `json:"role"`
	State    string `json:"state"`
}

type UserUpdateDto struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"omitempty,email"`
	Lastname string `json:"lastname"`
	Role     string `json:"role"`
	State    string `json:"state"`
}

type UserChgPasswdDto struct {
	Passwd    string `json:"password" validate:"required"`
	NPasswd   string `json:"new_password" validate:"required"`
	ReNPasswd string `json:"re_new_password" validate:"required"`
}

type UserAddPermissionDto struct {
	UserUuid string `json:"user_uuid" validate:"required"`
	DirUuid  string `json:"dir_uuid" validate:"required"`
}
