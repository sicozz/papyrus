package dtos

type DirDuplicateDto struct {
	Uuid      string `json:"uuid" validate:"uuid"`
	ParentDir string `json:"parent_dir" validate:"uuid"`
}
