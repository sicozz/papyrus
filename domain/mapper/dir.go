package mapper

import (
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
)

// Transform a dir entity into a dir dto
func MapDirToDirGetDto(d domain.Dir, path string, nChild, depth int) dtos.DirGetDto {
	return dtos.DirGetDto{
		Uuid:      d.Uuid,
		Name:      d.Name,
		ParentDir: d.ParentDir,
		Path:      path,
		Nchild:    nChild,
		Depth:     depth,
		Type:      "carpeta",
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
