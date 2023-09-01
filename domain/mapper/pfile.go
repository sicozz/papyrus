package mapper

import (
	"strconv"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils/constants"
)

// Transform a pfile entity into a pfile dto
func MapPFileToPFileGetDto(pf domain.PFile, apps []domain.Approvation) dtos.PFileGetDto {
	ap1 := ""
	ap2 := ""
	ap3 := ""
	chk1 := ""
	chk2 := ""
	chk3 := ""
	if len(apps) == 0 {
		return dtos.PFileGetDto{}
	}
	if len(apps) > 0 {
		ap1 = apps[0].UserUuid
		chk1 = strconv.FormatBool(apps[0].IsApproved)
	}
	if len(apps) > 1 {
		ap2 = apps[1].UserUuid
		chk2 = strconv.FormatBool(apps[1].IsApproved)
	}
	if len(apps) > 2 {
		ap3 = apps[2].UserUuid
		chk3 = strconv.FormatBool(apps[2].IsApproved)
	}

	return dtos.PFileGetDto{
		Uuid:         pf.Uuid,
		Code:         pf.Code,
		Name:         pf.Name,
		FsPath:       pf.FsPath,
		DateCreation: pf.DateCreation.Format(constants.LayoutDate),
		DateInput:    pf.DateInput.Format(constants.LayoutDate),
		Type:         pf.Type,
		State:        pf.State,
		Dir:          pf.Dir,
		RespUser:     pf.RespUser,
		AppUser1:     ap1,
		AppUser2:     ap2,
		AppUser3:     ap3,

		Chk1:    chk1,
		Chk2:    chk2,
		Chk3:    chk3,
		Version: pf.Version,
		Term:    pf.Term,
		Subtype: pf.Subtype,
	}
}

func MapPFileToDir(pf domain.PFile) domain.Dir {
	return domain.Dir{
		Uuid:      pf.Uuid,
		Name:      pf.Name,
		ParentDir: pf.Dir,
	}
}
