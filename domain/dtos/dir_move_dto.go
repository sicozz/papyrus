package dtos

type DirMoveDto struct {
	ParentDir string `json:"parent_dir" validate:"omitempty,uuid"`
}
