package domain

type ChgPasswd struct {
	Passwd    string `json:"password" validate:"required"`
	NPasswd   string `json:"new_password" validate:"required"`
	ReNPasswd string `json:"re_new_password" validate:"required"`
}

func NewChgPassw(passwd, nPasswd, reNPasswd string) (dto ChgPasswd) {
	dto = ChgPasswd{passwd, nPasswd, reNPasswd}
	return
}
