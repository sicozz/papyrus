package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/domain/mapper"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type pFileUseCase struct {
	pFileRepo      domain.PFileRepository
	dirRepo        domain.DirRepository
	userRepo       domain.UserRepository
	contextTimeout time.Duration
	log            utils.AggregatedLogger
}

// NewPFileUsecase will create a new dirUsecase object representation of domain.PFileUsecase interface
func NewPFileUsecase(pfr domain.PFileRepository, dr domain.DirRepository, ur domain.UserRepository, timeout time.Duration) domain.PFileUsecase {
	logger := utils.NewAggregatedLogger(constants.Usecase, constants.PFile)
	return &pFileUseCase{
		pFileRepo:      pfr,
		dirRepo:        dr,
		userRepo:       ur,
		contextTimeout: timeout,
		log:            logger,
	}
}

func (u *pFileUseCase) GetAll(c context.Context) (res []dtos.PFileGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	pFiles, err := u.pFileRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get dirs ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	pFilesDtos := make([]dtos.PFileGetDto, len(pFiles), len(pFiles))
	for i, pf := range pFiles {
		apps, err := u.pFileRepo.GetApprovations(ctx, pf.Uuid)
		if err != nil {
			u.log.Err("IN [GetAll] failed to get file approvations ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		}

		pFilesDtos[i] = mapper.MapPFileToPFileGetDto(pf, apps)
	}

	res = pFilesDtos

	return
}

func (u *pFileUseCase) Upload(c context.Context, p dtos.PFileUploadDto, file *multipart.FileHeader) (dto dtos.PFileGetDto, rErr domain.RequestErr) {
	// TODO: Refactor functions into something cleaner
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	src, err := file.Open()
	if err != nil {
		u.log.Err("IN [Upload] failed to open file ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	nFileUuid := uuid.New().String()
	nFilepath := constants.PathFsDir + string(os.PathSeparator) + nFileUuid + constants.UuidFileSeparator + file.Filename

	dst, err := os.Create(nFilepath)
	if err != nil {
		u.log.Err("IN [Upload] failed to create file ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		u.log.Err("IN [Upload] failed to copy file contents ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
	}

	// parse date_creation
	dateCreation, err := time.Parse(constants.LayoutDate, p.DateCreation)
	if err != nil {
		u.log.Err("IN [Upload] failed to parse DateCreation ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	// check type
	if exists := u.pFileRepo.ExistsTypeByDesc(ctx, p.Type); !exists {
		err := errors.New(fmt.Sprint("File type not found. type: ", p.Type))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	// check dir
	if exists := u.dirRepo.ExistsByUuid(ctx, p.Dir); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", p.Dir))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	// check responsible user
	if exists := u.userRepo.ExistsByUuid(ctx, p.RespUser); !exists {
		err := errors.New(fmt.Sprint("Responsible user not found. uuid: ", p.RespUser))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	approvations := []domain.Approvation{}
	if p.AppUser1 != "" {
		approvation := domain.Approvation{
			UserUuid:   p.AppUser1,
			IsApproved: p.Chk1,
		}
		approvations = append(approvations, approvation)
	}
	if p.AppUser2 != "" {
		approvation := domain.Approvation{
			UserUuid:   p.AppUser2,
			IsApproved: p.Chk2,
		}
		approvations = append(approvations, approvation)
	}
	if p.AppUser3 != "" {
		approvation := domain.Approvation{
			UserUuid:   p.AppUser3,
			IsApproved: p.Chk3,
		}
		approvations = append(approvations, approvation)
	}
	if len(approvations) == 0 {
		err := errors.New("File needs at least one approval user")
		rErr = domain.NewUCaseErr(http.StatusBadRequest, err)
		return
	}

	// check approval users
	for i, ap := range approvations {
		if exists := u.userRepo.ExistsByUuid(ctx, ap.UserUuid); !exists {
			err := errors.New(fmt.Sprint("Approval user ", i, " not found. uuid: ", ap.UserUuid))
			rErr = domain.NewUCaseErr(http.StatusNotFound, err)
			return
		}
	}

	// check name not taken in dir
	if taken := u.pFileRepo.IsNameTaken(ctx, p.Name, p.Dir); taken {
		err := errors.New("File name already taken")
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	// check code not taken
	if taken := u.pFileRepo.ExistsByCode(ctx, p.Code); taken {
		err := errors.New("File code already taken")
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	nPFile := domain.PFile{
		Uuid:         nFileUuid,
		Code:         p.Code,
		Name:         p.Name,
		RespUser:     p.RespUser,
		FsPath:       nFilepath,
		DateCreation: dateCreation,
		DateInput:    time.Now(),
		Type:         "TODO",
		State:        "TODO",
		Dir:          p.Dir,
		Version:      p.Version,
		Term:         p.Term,
		Subtype:      p.Subtype,
	}
	for i := range approvations {
		approvations[i].PFileUuid = nFileUuid
	}

	nUuid, err := u.pFileRepo.StoreUuid(ctx, nPFile, approvations)
	if err != nil {
		u.log.Err("IN [Upload] failed store pfile ->", err)
		err := errors.New("Could not upload file")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)

		// Delete created file
		err = os.Remove(nFilepath)
		if err != nil {
			u.log.Err("IN [Upload] failed remove created file ->", err)
			u.log.Wrn("IN [Upload] bad state. File created with no db representation ->", err)
			err := errors.New("Could not upload file")
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		}
		return
	}

	// WARN: Change this to usecase getbyuuid
	dto, err = u.GetByUuid(ctx, nUuid)
	if err != nil {
		u.log.Err("IN [Upload] failed fetch pfile ->", err)
		err := errors.New("Could not fetch file")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *pFileUseCase) GetByUuid(c context.Context, uuid string) (pFDto dtos.PFileGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	pFile, err := u.pFileRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetByUuid] failed to fetch pfile {", uuid, "} ->", err)
		err = errors.New("File not found. uuid: " + uuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	apps, err := u.pFileRepo.GetApprovations(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetByUuid] failed to get file approvations ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
	}

	pFDto = mapper.MapPFileToPFileGetDto(pFile, apps)

	return
}

func (u *pFileUseCase) Delete(c context.Context, uuid string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// TODO: Add checks before delete

	err := u.pFileRepo.Delete(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Delete] failed to delete pfile {", uuid, "} ->", err)
		err = errors.New("Failed to delete file")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *pFileUseCase) Approve(c context.Context, pfUuid, userUuid string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// TODO: Add exists check
	_, err := u.pFileRepo.GetByUuid(ctx, pfUuid)
	if err != nil {
		u.log.Err("IN [Approve] failed to fetch pfile {", pfUuid, "} ->", err)
		err = errors.New("File not found. uuid: " + pfUuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if exists := u.pFileRepo.ApprExistsByPK(ctx, pfUuid, userUuid); !exists {
		err = errors.New(fmt.Sprintf("User %v is not an approvation user of file %v", userUuid, pfUuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	err = u.pFileRepo.Approve(ctx, pfUuid, userUuid)
	if err != nil {
		u.log.Err("IN [Approve] failed to approve pfile ", pfUuid, " with user ", userUuid, " -> ", err)
		err = errors.New("Failed to delete file")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *pFileUseCase) Activate(c context.Context, pfUuid, userUuid string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistsByUuid(ctx, userUuid); !exists {
		err := errors.New("Responsible user not found. uuid: " + userUuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	// TODO: Add exists check
	_, err := u.pFileRepo.GetByUuid(ctx, pfUuid)
	if err != nil {
		u.log.Err("IN [Activate] failed to fetch pfile {", pfUuid, "} ->", err)
		err = errors.New("File not found. uuid: " + pfUuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if approved := u.pFileRepo.IsApproved(ctx, pfUuid); !approved {
		err = errors.New("File has not been approved")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	err = u.pFileRepo.Activate(ctx, pfUuid, userUuid)
	if err != nil {
		u.log.Err("IN [Activate] failed to activate pfile ", pfUuid, " -> ", err)
		err = errors.New("Failed to delete file")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}
