package dtos

type PFileGetDto struct {
	Uuid         string `json:"uuid"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	FsPath       string `json:"fs_path"`
	DateCreation string `json:"date_create"`
	DateInput    string `json:"date_input"`
	Type         string `json:"type"`
	State        string `json:"state"`
	Dir          string `json:"dir"`
	RespUser     string `json:"responsible_user"`
	AppUser1     string `json:"approval_user1"`
	AppUser2     string `json:"approval_user2"`
	AppUser3     string `json:"approval_user3"`

	// TODO: make this a checks of boolean type
	Chk1    string `json:"user_check1"`
	Chk2    string `json:"user_check2"`
	Chk3    string `json:"user_check3"`
	Version string `json:"version"`
	Term    int    `json:"term"`
	Subtype string `json:"subtype"`
}

type PFileUploadDto struct {
	Code         string `json:"code" validate:"required,ascii"`
	Name         string `json:"name" validate:"required,ascii"`
	DateCreation string `json:"date_create" validate:"required,datetime=2006-01-02"`
	Type         string `json:"type" validate:"required,ascii"`
	Dir          string `json:"dir" validate:"required,uuid"`
	RespUser     string `json:"responsible_user" validate:"required,uuid"`
	AppUser1     string `json:"approval_user1" validate:"required,uuid"`
	AppUser2     string `json:"approval_user2" validate:"omitempty,uuid"`
	AppUser3     string `json:"approval_user3" validate:"omitempty,uuid"`

	Chk1    bool   `json:"user_check1"`
	Chk2    bool   `json:"user_check2"`
	Chk3    bool   `json:"user_check3"`
	Version string `json:"version" validate:"required,ascii"`
	Term    int    `json:"term" validate:"required,number"`
	Subtype string `json:"subtype" validate:"required,ascii"`
}

type PFileChgCheckDto struct {
	Chk bool `json:"checked" validate:"required,boolean"`
}

type PFileChgStateDto struct {
	StateDesc string `json:"state" validate:"required,ascii"`
}
