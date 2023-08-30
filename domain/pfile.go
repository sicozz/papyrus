package domain

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/sicozz/papyrus/domain/dtos"
)

// PFile represents the File data struct
type PFile struct {
	Uuid         string
	Code         string
	Name         string
	FsPath       string // path type?
	DateCreation time.Time
	DateInput    time.Time
	Type         string
	State        string
	Stage        string
	Dir          string
	RevUser      string // revision user
	AppUser      string // approval user
}

type PFileUsecase interface {
	GetAll(c context.Context) ([]dtos.PFileGetDto, RequestErr)
	// GetByUuid(c context.Context, uuid string) (dtos.PFileGetDto, RequestErr)
	// Update(c context.Context, uuid string, p dtos.PFileUpdateDto) RequestErr
	Delete(c context.Context, uuid string) RequestErr

	Upload(c context.Context, p dtos.PFileUploadDto, file *multipart.FileHeader) (dtos.PFileGetDto, RequestErr)
	GetByUuid(ctx context.Context, uuid string) (PFile, RequestErr)
}

type PFileRepository interface {
	GetAll(ctx context.Context) ([]PFile, error)
	GetByUuid(ctx context.Context, uuid string) (PFile, error)
	// Update(ctx context.Context, uuid string, p dtos.PFileUpdateDto) error
	Delete(ctx context.Context, uuid string) error
	Store(ctx context.Context, pf PFile) (string, error)
	StoreUuid(ctx context.Context, pf PFile) (string, error)

	ExistsByCode(ctx context.Context, code string) bool
	IsNameTaken(ctx context.Context, name string, dirUuid string) bool

	// pfile_type ops
	ExistsTypeByDesc(ctx context.Context, desc string) bool

	// pfile_state ops
	ExistsStateByDesc(ctx context.Context, desc string) bool

	// pfile_stage ops
	ExistsStageByDesc(ctx context.Context, desc string) bool
}
