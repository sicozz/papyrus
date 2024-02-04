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
	Dir          string
	Version      string
	Term         int
	Subtype      string
	RespUser     string
}

type Approvation struct {
	UserUuid   string
	PFileUuid  string
	IsApproved bool
}

type Evidence struct {
	TaskUuid     string
	PFileUuid    string
	PFileName    string
	PFileFsPath  string
	DateCreation time.Time
}

type PFileUsecase interface {
	GetAll(c context.Context) ([]dtos.PFileGetDto, RequestErr)
	// GetByUuid(c context.Context, uuid string) (dtos.PFileGetDto, RequestErr)
	// Update(c context.Context, uuid string, p dtos.PFileUpdateDto) RequestErr
	Delete(c context.Context, uuid string) RequestErr

	Upload(c context.Context, p dtos.PFileUploadDto, file *multipart.FileHeader) (dtos.PFileGetDto, RequestErr)
	GetByUuid(c context.Context, uuid string) (dtos.PFileGetDto, RequestErr)

	ChgApprovation(c context.Context, pfUuid, userUuid string, chk bool) RequestErr
	ChgState(c context.Context, pfUuid, userUuid, stateDesc string) RequestErr

	RequestDownload(c context.Context, pfUuid, userUuid string) (string, RequestErr)
	AddDwnHistory(c context.Context, pfUuid, userUuid string) RequestErr

	GetEvidence(c context.Context, tUuid string) ([]dtos.PFileGetEvidenceDto, RequestErr)
	UploadEvidence(c context.Context, tUuid string, p dtos.PFileUploadDto, file *multipart.FileHeader) (dtos.PFileGetDto, RequestErr)
	DeleteEvidence(c context.Context, tUuid, pfUuid string) RequestErr
}

type PFileRepository interface {
	GetAll(ctx context.Context) ([]PFile, error)
	GetByUuid(ctx context.Context, uuid string) (PFile, error)
	GetByUser(ctx context.Context, uuid string) ([]PFile, error)
	GetByDir(ctx context.Context, uuid string) ([]PFile, error)

	// Update(ctx context.Context, uuid string, p dtos.PFileUpdateDto) error
	Delete(ctx context.Context, uuid string) error
	StoreUuid(ctx context.Context, pf PFile, apps []Approvation) (string, error)

	GetApprovations(c context.Context, uuid string) ([]Approvation, error)
	ChgApprovation(ctx context.Context, pfUuid, userUuid string, chk bool) error
	ChgState(ctx context.Context, pfUuid, userUuid, stateDesc string) error
	ChgStateBypass(ctx context.Context, pfUuid, stateDesc string) error
	ApprExistsByPK(ctx context.Context, pfUuid, userUuid string) bool

	ExistsByUuid(ctx context.Context, uuid string) bool
	ExistsByCode(ctx context.Context, code string) bool
	IsNameTaken(ctx context.Context, name string, dirUuid string) bool
	IsApproved(ctx context.Context, uuid string) bool

	AddDwnHistory(ctx context.Context, date time.Time, pfUuid, userUuid string) error

	GetEvidence(ctx context.Context, tUuid string) ([]Evidence, error)
	AddEvidence(ctx context.Context, tUuid, pfUuid string) error
	DeleteEvidence(ctx context.Context, tUuid, pfUuid string) error

	// pfile_type ops
	ExistsTypeByDesc(ctx context.Context, desc string) bool

	// pfile_state ops
	ExistsStateByDesc(ctx context.Context, desc string) bool

	// pfile_stage ops
	ExistsStageByDesc(ctx context.Context, desc string) bool
}
