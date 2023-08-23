package mapper

import (
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
)

// Transform a dir entity into a dir dto
func MapDirToDirGetDto(d domain.Dir) dtos.DirGetDto {
	t := dtos.DirGetDtoType
	if d.Name[0] == '_' {
		t = "documento"
		d.Name = d.Name[1:]
	}
	return dtos.DirGetDto{
		Uuid:      d.Uuid,
		Name:      d.Name,
		ParentDir: d.ParentDir,
		Path:      d.Path,
		Nchild:    d.Nchild,
		Depth:     d.Depth,
		Type:      t,
		Visible:   dtos.DirGetDtoVisible,
		Open:      dtos.DirGetDtoOpen,
	}
}

// Transform a DirStoreDto entity into a dir
func MapDirStoreDtoToDir(p dtos.DirStoreDto) domain.Dir {
	return domain.Dir{
		Name:      p.Name,
		ParentDir: p.ParentDir,
	}
}
