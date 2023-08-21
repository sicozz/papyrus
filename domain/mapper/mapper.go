package mapper

import (
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
)

// Transform a dir entity into a dir dto
func MapDirToDirGetDto(d domain.Dir) dtos.DirGetDto {
	return dtos.DirGetDto{
		Uuid:      d.Uuid,
		Name:      d.Name,
		ParentDir: d.ParentDir,
		Path:      d.Path,
		Nchild:    d.Nchild,
		Depth:     d.Depth,
		Type:      dtos.DirGetDtoType,
		Visible:   dtos.DirGetDtoVisible,
		Open:      dtos.DirGetDtoOpen,
	}
}

// Transform a dir entity into a dir dto
func MapUserToUserGetDto(u domain.User) dtos.UserGetDto {
	return dtos.UserGetDto{
		Uuid:     u.Uuid,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Lastname: u.Lastname,
		Role:     u.Role.Description,
		State:    u.State.Description,
	}
}
