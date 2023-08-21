package dtos

type ChgPasswdDto struct {
	Passwd    string `json:"password" validate:"required"`
	NPasswd   string `json:"new_password" validate:"required"`
	ReNPasswd string `json:"re_new_password" validate:"required"`
}

func NewChgPassw(passwd, nPasswd, reNPasswd string) (dto ChgPasswdDto) {
	dto = ChgPasswdDto{passwd, nPasswd, reNPasswd}
	return
}
