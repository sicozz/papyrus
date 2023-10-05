package mapper

import (
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils/constants"
)

// Transform a dir entity into a dir dto
func MapUserToUserGetDto(u domain.User) dtos.UserGetDto {
	return dtos.UserGetDto{
		Uuid:     u.Uuid,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Lastname: u.Lastname,
		Role: dtos.RoleDto{
			Code:        u.Role.Code,
			Description: u.Role.Description,
		},
		State: dtos.UStateDto{
			Code:        u.State.Code,
			Description: u.State.Description,
		},
	}
}

// Transform a dir entity into a dir dto
func MapUserStoreDtoToUser(u dtos.UserStore) domain.User {
	return domain.User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Lastname: u.Lastname,
		Role:     domain.Role{},
		State:    domain.UserState{},
	}
}

func MapHistoryToUserHistoryGetDto(h domain.History) dtos.UserHistoryGetDto {
	return dtos.UserHistoryGetDto{
		Date:      h.Date.Format(constants.LayoutDate),
		UserUuid:  h.UserUuid,
		PFileUuid: h.PFileUuid,
	}
}
