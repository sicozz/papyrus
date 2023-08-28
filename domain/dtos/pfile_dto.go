package dtos

type PFileGetDto struct {
	Uuid         string `json:"uuid"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	DateCreation string `json:"date_create"`
	DateInput    string `json:"date_input"`
	Type         string `json:"type"`
	State        string `json:"state"`
	Stage        string `json:"stage"`
	Dir          string `json:"dir"`
	RevUser      string `json:"responsible_user"` // revision user
	AppUser      string `json:"approval_user"`    // approval user
}

type PFileUploadDto struct {
	Code         string `json:"code" validate:"required,ascii"`
	Name         string `json:"name" validate:"required,ascii"`
	DateCreation string `json:"date_create" validate:"required,datetime=2006-01-02"`
	Type         string `json:"type" validate:"required,ascii"`
	Dir          string `json:"dir" validate:"required,uuid"`
	RevUser      string `json:"responsible_user" validate:"required,uuid"` // revision user
	AppUser      string `json:"approval_user" validate:"required,uuid"`    // approval user
}
