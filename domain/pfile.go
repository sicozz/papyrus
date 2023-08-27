package domain

import (
	"context"
	"time"

	"github.com/sicozz/papyrus/domain/dtos"
)

// PFile represents the File data struct
type PFile struct {
	Uuid         string
	Code         string
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
	// Delete(c context.Context, uuid string) RequestErr

	// Upload(c context.Context, d dtos.PFileUploadDto) (dtos.PFileGetDto, RequestErr)
	// Download(c context.Context, uuid string) (dtos.PFileGetDto, RequestErr)
}

type PFileRepository interface {
	GetAll(c context.Context) ([]PFile, error)
	// GetByUuid(c context.Context, uuid string) (PFile, error)
	// Update(c context.Context, uuid string, p dtos.PFileUpdateDto) error
	// Delete(c context.Context, uuid string) error

	// Upload(c context.Context, d dtos.PFileUploadDto) (PFile, error)
	// Download(c context.Context, uuid string) (PFile, error)
}
