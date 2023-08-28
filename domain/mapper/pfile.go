package mapper

import (
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils/constants"
)

// Transform a pfile entity into a pfile dto
func MapPFileToPFileGetDto(pf domain.PFile) dtos.PFileGetDto {
	return dtos.PFileGetDto{
		Uuid:         pf.Uuid,
		Code:         pf.Code,
		Name:         pf.Name,
		DateCreation: pf.DateCreation.Format(constants.LayoutDate),
		DateInput:    pf.DateInput.Format(constants.LayoutDate),
		Type:         pf.Type,
		State:        pf.State,
		Stage:        pf.Stage,
		Dir:          pf.Dir,
		RevUser:      pf.RevUser,
		AppUser:      pf.AppUser,
	}
}

func MapPFileToDir(pf domain.PFile) domain.Dir {
	return domain.Dir{
		Uuid:      pf.Uuid,
		Name:      pf.Name,
		ParentDir: pf.Dir,
	}
}
