package dtos

type TaskGetDto struct {
	Uuid         string `json:"uuid"`
	Name         string `json:"name"`
	Procedure    string `json:"procedure"`
	DateCreation string `json:"date_create"`
	DateCheck    string `json:"date_check"`
	DateClose    string `json:"date_close"`
	Term         int    `json:"term"`
	State        string `json:"state"`
	Dir          string `json:"dir"`
	CreatorUser  string `json:"creator_user"`
	RecvUser     string `json:"receiver_user"`
	Chk          bool   `json:"check"`
	Plan         string `json:"plan"`
}

type TaskStoreDto struct {
	Name        string `json:"name" validate:"required,excludesall=!@?"`
	Procedure   string `json:"procedure" validate:"required,excludesall=!@?"`
	DateCreate  string `json:"date_create" validate:"required,datetime=2006-01-02"`
	Term        int    `json:"term" validate:"required,number"`
	Dir         string `json:"dir" validate:"required,uuid"`
	CreatorUser string `json:"creator_user" validate:"required,uuid"`
	RecvUser    string `json:"receiver_user" validate:"required,uuid"`
	Plan        string `json:"plan"`
}

type TaskChgCheck struct {
	Chk bool `json:"checked" validate:"boolean"`
}

type TaskChgStateDto struct {
	StateDesc string `json:"state" validate:"required,excludesall=!@#?"`
}

type Tmptest struct {
	Name string `json:"name"`
}
