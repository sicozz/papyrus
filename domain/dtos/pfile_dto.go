package dtos

type PFileGetDto struct {
	Uuid         string `json:"uuid"`
	Code         string `json:"code"`
	DateCreation string `json:"date_create"`
	DateInput    string `json:"date_input"`
	Type         string `json:"type"`
	State        string `json:"state"`
	Stage        string `json:"stage"`
	Dir          string `json:"dir"`
	RevUser      string `json:"responsible_user"` // revision user
	AppUser      string `json:"approval_user"`    // approval user
}
