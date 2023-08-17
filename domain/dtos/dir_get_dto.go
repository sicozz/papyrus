package dtos

import "github.com/sicozz/papyrus/domain"

const (
	dirGetDtoType    = "carpeta"
	dirGetDtoVisible = false
	dirGetDtoOpen    = false
)

// DirGetDto represents the information served in the GetAll endpoint of DirHandler
type DirGetDto struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name" validate:"required,ascii"`
	ParentDir string `json:"parent_dir" validate:"required,ascii,uuid"`
	Path      string `json:"path"`
	Nchild    int    `json:"nchild"`
	Depth     int    `json:"depth"`
	Type      string `json:"type"`
	Visible   bool   `json:"visible"`
	Open      bool   `json:"open"`
}

func NewDirGetDto(dir domain.Dir) DirGetDto {
	return DirGetDto{
		Uuid:      dir.Uuid,
		Name:      dir.Name,
		ParentDir: dir.ParentDir,
		Path:      dir.Path,
		Nchild:    dir.Nchild,
		Depth:     dir.Depth,
		Type:      dirGetDtoType,
		Visible:   dirGetDtoVisible,
		Open:      dirGetDtoOpen,
	}
}
