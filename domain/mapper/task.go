package mapper

import (
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils/constants"
)

func MapTaskToTaskGetDto(t domain.Task) dtos.TaskGetDto {
	return dtos.TaskGetDto{
		Uuid:         t.Uuid,
		Name:         t.Name,
		Procedure:    t.Procedure,
		DateCreation: t.DateCreation.Format(constants.LayoutDate),
		Term:         t.Term,
		State:        t.State,
		Dir:          t.Dir,
		CreatorUser:  t.CreatorUser,
		RecvUser:     t.RecvUser,
		Chk:          t.Check,
	}
}

func MapTaskToDir(t domain.Task) domain.Dir {
	return domain.Dir{
		Uuid:      t.Uuid,
		Name:      t.Name,
		ParentDir: t.Dir,
	}
}

func MapTaskToDocsNotDirGetDto(t domain.Task) dtos.DocsNotDirGetDto {
	return dtos.DocsNotDirGetDto{
		Uuid:         t.Uuid,
		Name:         t.Name,
		Procedure:    t.Procedure,
		DateCreation: t.DateCreation.Format(constants.LayoutDate),
		Term:         t.Term,
		State:        t.State,
		ParentDir:    t.Dir,
		CreatorUser:  t.CreatorUser,
		RecvUser:     t.RecvUser,
		Chk:          t.Check,
	}
}
