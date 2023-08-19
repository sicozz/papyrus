package instances

import "github.com/sicozz/papyrus/domain"

func GetUser() domain.User {
	return domain.User{
		Uuid:     "00000000-0000-0000-0000-000000000000",
		Username: "tron",
		Email:    "tron@emcom.com",
		Password: "0pen6ate",
		Name:     "Sam",
		Lastname: "Flynn",
		Role:     domain.Role{Code: 3},
		State:    domain.UserState{Code: 2},
	}
}

func GetUserList() []domain.User {
	return []domain.User{
		{
			Uuid:     "11111111-1111-1111-1111-111111111111",
			Username: "doom",
			Email:    "doom@email.com",
			Password: "iconofsin",
			Name:     "John",
			Lastname: "Carmack",
			Role:     domain.Role{Code: 3},
			State:    domain.UserState{Code: 2},
		},
		{
			Uuid:     "22222222-2222-2222-2222-222222222222",
			Username: "vi",
			Email:    "bj@berkeley.com",
			Password: "ipasswdesc",
			Name:     "Bill",
			Lastname: "Joy",
			Role:     domain.Role{Code: 2},
			State:    domain.UserState{Code: 1},
		},
		{
			Uuid:     "33333333-3333-3333-3333-333333333333",
			Username: "linux",
			Email:    "linus@unix.com",
			Password: "tux",
			Name:     "Linus",
			Lastname: "Torvalds",
			Role:     domain.Role{Code: 1},
			State:    domain.UserState{Code: 2},
		},
	}
}

func GetRole() domain.Role {
	return domain.Role{Code: 2, Description: "admin"}
}

func GetRoleList() []domain.Role {
	return []domain.Role{
		{Code: 1, Description: "estandar"},
		{Code: 2, Description: "admin"},
		{Code: 3, Description: "super"},
	}
}

func GetUState() domain.UserState {
	return domain.UserState{Code: 2, Description: "admin"}
}

func GetUStateList() []domain.UserState {
	return []domain.UserState{
		{Code: 1, Description: "estandar"},
		{Code: 2, Description: "admin"},
		{Code: 3, Description: "super"},
	}
}
