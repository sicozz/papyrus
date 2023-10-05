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

type UserHistoryGetDto struct {
	DownloadUuid   string `json:"download_uuid"`
	Date           string `json:"download_date"`
	UserUuid       string `json:"user_uuid"`
	PFileUuid      string `json:"file_uuid"`
	PFileCode      string `json:"file_code"`
	PFileVersion   string `json:"file_version"`
	PFileTerm      string `json:"file_term"`
	PFileName      string `json:"file_name"`
	PFileType      string `json:"file_type"`
	PFileDateInput string `json:"file_date_input"`
	PFileDir       string `json:"file_dir_uuid"`
	PFileDirPath   string `json:"file_dir_path"`
}
