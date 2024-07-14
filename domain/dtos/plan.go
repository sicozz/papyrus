package dtos

type PlanGetDto struct {
	Uuid         string       `json:"uuid"`
	Code         string       `json:"code" validate:"required,excludesall=!@?"`
	Name         string       `json:"name" validate:"required,excludesall=!@?"`
	Origin       string       `json:"origin" validate:"required,excludesall=!@?"`
	ActionType   string       `json:"action_type" validate:"required,excludesall=!@?"`
	Term         int          `json:"term" validate:"required,excludesall=!@?"`
	CreatorUser  string       `json:"creator_user" validate:"required,excludesall=!@?"`
	RespUser     string       `json:"responsible_user" validate:"required,excludesall=!@?"`
	DateCreation string       `json:"date_create" validate:"required,excludesall=!@?"`
	DateCheck    string       `json:"date_check" validate:"required,excludesall=!@?"`
	DateClose    string       `json:"date_close" validate:"required,excludesall=!@?"`
	Causes       string       `json:"causes" validate:"required,excludesall=!@?"`
	Conclusions  string       `json:"conclusions" validate:"required,excludesall=!@?"`
	State        string       `json:"state" validate:"required,excludesall=!@?"`
	Stage        int          `json:"stage" validate:"required,excludesall=!@?"`
	Dir          string       `json:"parent_dir" validate:"required,excludesall=!@?"`
	Tasks        []TaskGetDto `json:"tasks"`

	Action0desc string `json:"action0_desc"`
	Action0date string `json:"action0_date"`
	Action0user string `json:"action0_user"`
	Action1desc string `json:"action1_desc"`
	Action1date string `json:"action1_date"`
	Action1user string `json:"action1_user"`
	Action2desc string `json:"action2_desc"`
	Action2date string `json:"action2_date"`
	Action2user string `json:"action2_user"`
	Action3desc string `json:"action3_desc"`
	Action3date string `json:"action3_date"`
	Action3user string `json:"action3_user"`
	Action4desc string `json:"action4_desc"`
	Action4date string `json:"action4_date"`
	Action4user string `json:"action4_user"`
	Action5desc string `json:"action5_desc"`
	Action5date string `json:"action5_date"`
	Action5user string `json:"action5_user"`
}

type PlanStoreDto struct {
	Code         string `json:"code" validate:"required,excludesall=!@?"`
	Name         string `json:"name" validate:"required,excludesall=!@?"`
	Origin       string `json:"origin" validate:"required,excludesall=!@?"`
	ActionType   string `json:"action_type" validate:"required,excludesall=!@?"`
	Term         int    `json:"term" validate:"required,excludesall=!@?"`
	CreatorUser  string `json:"creator_user" validate:"required,excludesall=!@?"`
	RespUser     string `json:"responsible_user" validate:"required,excludesall=!@?"`
	DateCreation string `json:"date_create" validate:"required,excludesall=!@?"`
	Dir          string `json:"parent_dir" validate:"required,excludesall=!@?"`
	State        string `json:"state" validate:"required,excludesall=!@?"`
	Stage        int    `json:"stage" validate:"required,excludesall=!@?"`
	Causes       string `json:"causes" validate:"excludesall=!@?"`
	Conclusions  string `json:"conclusions" validate:"excludesall=!@?"`

	Action0desc string `json:"action0_desc" validate:"excludesall=!@?"`
	Action0date string `json:"action0_date" validate:"excludesall=!@?"`
	Action0user string `json:"action0_user" validate:"excludesall=!@?"`
	Action1desc string `json:"action1_desc" validate:"excludesall=!@?"`
	Action1date string `json:"action1_date" validate:"excludesall=!@?"`
	Action1user string `json:"action1_user" validate:"excludesall=!@?"`
	Action2desc string `json:"action2_desc" validate:"excludesall=!@?"`
	Action2date string `json:"action2_date" validate:"excludesall=!@?"`
	Action2user string `json:"action2_user" validate:"excludesall=!@?"`
	Action3desc string `json:"action3_desc" validate:"excludesall=!@?"`
	Action3date string `json:"action3_date" validate:"excludesall=!@?"`
	Action3user string `json:"action3_user" validate:"excludesall=!@?"`
	Action4desc string `json:"action4_desc" validate:"excludesall=!@?"`
	Action4date string `json:"action4_date" validate:"excludesall=!@?"`
	Action4user string `json:"action4_user" validate:"excludesall=!@?"`
	Action5desc string `json:"action5_desc" validate:"excludesall=!@?"`
	Action5date string `json:"action5_date" validate:"excludesall=!@?"`
	Action5user string `json:"action5_user" validate:"excludesall=!@?"`
}

type PlanUpdateDto struct {
	Name        string `json:"name" validate:"required,excludesall=!@?"`
	Origin      string `json:"origin" validate:"required,excludesall=!@?"`
	ActionType  string `json:"action_type" validate:"required,excludesall=!@?"`
	Term        int    `json:"term" validate:"required,excludesall=!@?"`
	RespUser    string `json:"responsible_user" validate:"required,excludesall=!@?"`
	DateClose   string `json:"date_close" validate:"excludesall=!@?"`
	Causes      string `json:"causes" validate:"excludesall=!@?"`
	Conclusions string `json:"conclusions" validate:"excludesall=!@?"`
	Dir         string `json:"parent_dir" validate:"required,excludesall=!@?"`
	State       string `json:"state" validate:"required,excludesall=!@?"`
	Stage       int    `json:"stage" validate:"required,excludesall=!@?"`

	Action0desc string `json:"action0_desc" validate:"excludesall=!@?"`
	Action0date string `json:"action0_date" validate:"excludesall=!@?"`
	Action0user string `json:"action0_user" validate:"excludesall=!@?"`
	Action1desc string `json:"action1_desc" validate:"excludesall=!@?"`
	Action1date string `json:"action1_date" validate:"excludesall=!@?"`
	Action1user string `json:"action1_user" validate:"excludesall=!@?"`
	Action2desc string `json:"action2_desc" validate:"excludesall=!@?"`
	Action2date string `json:"action2_date" validate:"excludesall=!@?"`
	Action2user string `json:"action2_user" validate:"excludesall=!@?"`
	Action3desc string `json:"action3_desc" validate:"excludesall=!@?"`
	Action3date string `json:"action3_date" validate:"excludesall=!@?"`
	Action3user string `json:"action3_user" validate:"excludesall=!@?"`
	Action4desc string `json:"action4_desc" validate:"excludesall=!@?"`
	Action4date string `json:"action4_date" validate:"excludesall=!@?"`
	Action4user string `json:"action4_user" validate:"excludesall=!@?"`
	Action5desc string `json:"action5_desc" validate:"excludesall=!@?"`
	Action5date string `json:"action5_date" validate:"excludesall=!@?"`
	Action5user string `json:"action5_user" validate:"excludesall=!@?"`
}
