package domain

import (
	"context"

	"github.com/sicozz/papyrus/domain/dtos"
)

// Dir represents the Directory data strict
type Dir struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name" validate:"required,ascii"`
	ParentDir string `json:"parent_dir" validate:"required,ascii,uuid"`
}

// DirUsecase represents the dir's usecases
type DirUsecase interface {
	GetAll(c context.Context) ([]dtos.DirGetDto, RequestErr)
	GetByUuid(c context.Context, uuid string) (dtos.DirGetDto, RequestErr)
	GetDocsByUser(c context.Context, uuid string) ([]dtos.DocsNotDirGetDto, RequestErr)
	// GetBranch(c context.Context, uuid string) ([]dtos.DirGetDto, RequestErr)
	Store(c context.Context, d dtos.DirStoreDto) (dtos.DirGetDto, RequestErr)
	Update(c context.Context, uuid string, p dtos.DirUpdateDto) RequestErr
	Delete(c context.Context, uuid string) RequestErr
	Move(c context.Context, uuid string, nPUuid string) RequestErr
	Duplicate(c context.Context, p dtos.DirDuplicateDto) ([]dtos.DirGetDto, RequestErr)
}

// DirRepository represents the dir's repository contract
type DirRepository interface {
	// Calc fields
	GetNChild(ctx context.Context, uuid string) (int, error)
	GetPath(ctx context.Context, uuid string) (string, error)
	GetDepth(ctx context.Context, uuid string) (int, error)

	// Check constraints
	IsNameTaken(ctx context.Context, name string, destUuid string) bool
	IsSubDir(ctx context.Context, uuid string, destUuid string) bool

	GetAll(ctx context.Context) ([]Dir, error)
	// GetBranch(c context.Context, uuid string) ([]dtos.DirGetDto, RequestErr)
	GetByUuid(ctx context.Context, uuid string) (Dir, error)
	ExistsByUuid(ctx context.Context, uuid string) bool
	Store(ctx context.Context, d *Dir) (string, error)
	InsertDirs(ctx context.Context, dirs []Dir) error
	Delete(ctx context.Context, uuid string) error
	// DeleteBranch(c context.Context, uuid string) ([]dtos.DirGetDto, RequestErr)
	ChgName(ctx context.Context, uuid string, nName string) error
	ChgParentDir(ctx context.Context, uuid string, nPUuid string) error
}
