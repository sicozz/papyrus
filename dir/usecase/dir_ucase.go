package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type dirUsecase struct {
	dirRepo        domain.DirRepository
	contextTimeout time.Duration
	log            utils.AggregatedLogger
}

// NewDirUsecase will create a new dirUsecase object representation of domain.DirUsecase interface
func NewDirUsecase(dr domain.DirRepository, timeout time.Duration) domain.DirUsecase {
	logger := utils.NewAggregatedLogger(constants.Usecase, constants.Dir)
	return &dirUsecase{
		dirRepo:        dr,
		contextTimeout: timeout,
		log:            logger,
	}
}

func (u *dirUsecase) GetAll(c context.Context) (res []domain.Dir, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.dirRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll]: could not get dirs ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) GetByUuid(c context.Context, uuid string) (res domain.Dir, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.dirRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetByUuid]: could not get dir ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) Store(c context.Context, dir *domain.Dir) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.dirRepo.ExistByUuid(ctx, dir.ParentDir); !exists {
		err := errors.New("Parent dir not found")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	parentDir, err := u.dirRepo.GetByUuid(ctx, dir.ParentDir)
	if err != nil {
		err := errors.New("Could not fetch parent dir")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	/*
		TODO: Validate name not taken in parent_dir
	*/
	dirs := []string{parentDir.Path, dir.Name}
	dir.Path = strings.Join(dirs, string(os.PathSeparator))
	dir.Nchild = 0
	dir.Depth = parentDir.Depth + 1

	/* TODO: WARN: Make the storage of the dir and increment of parent_dir.nchild a transaction */
	err = u.dirRepo.Store(ctx, dir)
	if err != nil {
		err = errors.New("Dir creation failed")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	err = u.dirRepo.IncNchild(ctx, parentDir.Uuid, 1)
	if err != nil {
		err = errors.New("Dir nchild increment failed")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) Update(c context.Context, uuid string, dUp *domain.Dir) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.dirRepo.ExistByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}
	/* TODO: Add validation validate unique name in parent_dir */
	/* WARN: VALIDATE THAT ROOT FOLDER IS NOT CHANGED */

	if dUp.Name == "" {
		return
	}

	err := u.dirRepo.ChgName(ctx, uuid, dUp.Name)
	if err != nil {
		u.log.Err("IN [Update]: could not change name ->", err)
		err = errors.New(fmt.Sprint("Dir patch failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) Delete(c context.Context, uuid string) (rErr domain.RequestErr) {
	/* WARN: VALIDATE THAT THE ROOT DIR IS NOT DELETED */
	/* TODO: WARN: Make the deletion of the dir and decrement of parent_dir.nchild a transaction */
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.dirRepo.ExistByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	dir, err := u.dirRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Delete]: could get dir {", uuid, "} ->", err)
		err := errors.New("Could not fetch dir")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	parentDir, err := u.dirRepo.GetByUuid(ctx, dir.ParentDir)
	if err != nil {
		u.log.Err("IN [Delete]: could get parent dir {", uuid, "} ->", err)
		err := errors.New("Could not fetch parent dir")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	err = u.dirRepo.Delete(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Delete]: could not delete dir {", uuid, "} ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	err = u.dirRepo.DecNchild(ctx, parentDir.Uuid, 1)
	if err != nil {
		err = errors.New("Dir nchild increment failed")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) Move(c context.Context, uuid string, nPUuid string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.dirRepo.ExistByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	dir, err := u.dirRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Move]: could get dir {", uuid, "} ->", err)
		err := errors.New("Could not fetch dir")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	if exists := u.dirRepo.ExistByUuid(ctx, nPUuid); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	/* TODO: Add validation validate unique name in parent_dir */
	/* WARN: VALIDATE THAT ROOT FOLDER IS NOT MOVED */
	nDir, err := u.dirRepo.GetByUuid(ctx, nPUuid)
	if err != nil {
		u.log.Err("IN [Move]: could get new parent dir {", uuid, "} ->", err)
		err := errors.New("Could not fetch dir")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	// TODO: REFACTOR: Unify the move operation as a db transaction
	err = u.dirRepo.ChgParentDir(ctx, uuid, nPUuid)
	if err != nil {
		u.log.Err("IN [Move]: could not change parent dir ->", err)
		err = errors.New(fmt.Sprint("Dir patch failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	dirs := []string{nDir.Path, dir.Name}
	nPath := strings.Join(dirs, string(os.PathSeparator))
	u.log.Wrn("PATH>>>\t", nPath)
	err = u.dirRepo.ChgPath(ctx, uuid, nPath)
	if err != nil {
		u.log.Err("IN [Move]: could not change dir path ->", err)
		err = errors.New(fmt.Sprint("Dir patch failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	err = u.dirRepo.ChgDepth(ctx, uuid, nDir.Depth+1)
	if err != nil {
		u.log.Err("IN [Move]: could not change dir path ->", err)
		err = errors.New(fmt.Sprint("Dir patch failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	err = u.dirRepo.DecNchild(ctx, dir.ParentDir, 1)
	if err != nil {
		u.log.Err("IN [Move]: could not decrease old parent dir nchild ->", err)
		err = errors.New(fmt.Sprint("Dir patch failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	err = u.dirRepo.IncNchild(ctx, nPUuid, 1)
	if err != nil {
		u.log.Err("IN [Move]: could not increase new parent dir nchild ->", err)
		err = errors.New(fmt.Sprint("Dir patch failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) Duplicate(c context.Context, uuid string, destUuid string) (res domain.Dir, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	dirs, err := u.dirRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [Duplicate] failed to get dirs ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	neoDir, dupDirs, err := domain.Duplicate(uuid, destUuid, dirs)
	u.dirRepo.IncNchild(ctx, destUuid, 1)
	for _, d := range dupDirs {
		u.dirRepo.Insert(ctx, *d)
	}

	res = *neoDir

	// neoDirs := []domain.Dir{}
	// for _, d := range dupDirs {
	// 	neoDirs = append(neoDirs, *d)
	// }

	// newDirs := append(dirs, neoDirs...)
	// u.log.Inf(fmt.Sprintf(">>>New fs: %+v", newDirs))

	return
}
