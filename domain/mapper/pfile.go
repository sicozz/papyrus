package mapper

import (
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
)

// Transform a pfile entity into a pfile dto
func MapPFileToPFileGetDto(pf domain.PFile) dtos.PFileGetDto {
	return dtos.PFileGetDto{
		Uuid:         pf.Uuid,
		Code:         pf.Code,
		DateCreation: pf.DateCreation.Format("2006-01-02"),
		DateInput:    pf.DateInput.Format("2006-01-02"),
		Type:         pf.Type,
		State:        pf.State,
		Stage:        pf.Stage,
		Dir:          pf.Dir,
		RevUser:      pf.RevUser,
		AppUser:      pf.AppUser,
	}
}
