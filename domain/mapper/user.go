package mapper

import (
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
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
