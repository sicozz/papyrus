package dtos

const (
	DirGetDtoType    = "carpeta"
	DirGetDtoVisible = false
	DirGetDtoOpen    = false
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
	State     string `json:"state"`
}

type DirStoreDto struct {
	Name      string `json:"name" validate:"required,ascii"`
	ParentDir string `json:"parent_dir" validate:"required,ascii,uuid"`
}

type DirUpdateDto struct {
	Name string `json:"name" validate:"required,ascii"`
}

type DirDuplicateDto struct {
	Uuid      string `json:"uuid" validate:"uuid"`
	ParentDir string `json:"parent_dir" validate:"uuid"`
	Name      string `json:"name" validate:"ascii"`
}

type DirMoveDto struct {
	ParentDir string `json:"parent_dir" validate:"omitempty,uuid"`
}
