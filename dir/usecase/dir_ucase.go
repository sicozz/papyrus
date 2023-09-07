package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/domain/mapper"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type dirUsecase struct {
	dirRepo        domain.DirRepository
	pFileRepo      domain.PFileRepository
	taskRepo       domain.TaskRepository
	contextTimeout time.Duration
	log            utils.AggregatedLogger
}

// NewDirUsecase will create a new dirUsecase object representation of domain.DirUsecase interface
func NewDirUsecase(dr domain.DirRepository, pfr domain.PFileRepository, tr domain.TaskRepository, timeout time.Duration) domain.DirUsecase {
	logger := utils.NewAggregatedLogger(constants.Usecase, constants.Dir)
	return &dirUsecase{
		dirRepo:        dr,
		pFileRepo:      pfr,
		taskRepo:       tr,
		contextTimeout: timeout,
		log:            logger,
	}
}

func (u *dirUsecase) GetAll(c context.Context) (res []dtos.DirGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	dirs, err := u.dirRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get dirs ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	pfiles, err := u.pFileRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get files ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}
	for _, pf := range pfiles {
		pFileDirDto := mapper.MapPFileToDir(pf)
		dirs = append(dirs, pFileDirDto)
	}

	tasks, err := u.taskRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get tasks ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}
	for _, t := range tasks {
		taskDirDto := mapper.MapTaskToDir(t)
		dirs = append(dirs, taskDirDto)
	}

	nChild, err := u.dirRepo.GetNChild(ctx, constants.RootDirUuid)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get root children number ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res, err = fillDetailsBFS(constants.RootDirUuid, dirs, "", nChild, 0)
	if err != nil {
		u.log.Err("IN [GetAll] failed to fill tree details ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	// TODO: Refactor folder and document mix
	for _, pf := range pfiles {
		for j := range res {
			if pf.Uuid == res[j].Uuid {
				// TODO: put this type and state in mapper
				res[j].Type = "documento"
				res[j].State = pf.State
			}
		}
	}
	for _, t := range tasks {
		for j := range res {
			if t.Uuid == res[j].Uuid {
				// TODO: put this type and state in mapper
				res[j].Type = "tarea"
				res[j].State = t.State
			}
		}
	}
	return
}

func (u *dirUsecase) GetByUuid(c context.Context, uuid string) (res dtos.DirGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	dir, err := u.dirRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetByUuid] failed to get dir ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	path, err := u.dirRepo.GetPath(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetByUuid] failed to get dir path ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	nChild, err := u.dirRepo.GetNChild(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetByUuid] failed to get dir children number ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	depth, err := u.dirRepo.GetDepth(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetByUuid] failed to get dir depth ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = mapper.MapDirToDirGetDto(dir, path, nChild, depth)

	return
}

func (u *dirUsecase) Store(c context.Context, p dtos.DirStoreDto) (res dtos.DirGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.dirRepo.ExistsByUuid(ctx, p.ParentDir); !exists {
		err := errors.New("Parent dir not found")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if taken := u.dirRepo.IsNameTaken(ctx, p.Name, p.ParentDir); taken {
		err := errors.New(fmt.Sprint("Name already taken in destination dir. name: ", p.Name))
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	dir := mapper.MapDirStoreDtoToDir(p)

	uuid, err := u.dirRepo.Store(ctx, &dir)
	if err != nil {
		err = errors.New("Dir creation failed")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res, err = u.GetByUuid(ctx, uuid)
	if err != nil {
		err = errors.New("Dir fetch failed")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) Update(c context.Context, uuid string, p dtos.DirUpdateDto) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if constants.RootDirUuid == uuid {
		err := errors.New("Is not possible to modify the root directory")
		rErr = domain.NewUCaseErr(http.StatusConflict, err)
		return
	}

	if exists := u.dirRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	dir, err := u.dirRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Update] failed to fetch dir ->", err)
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if taken := u.dirRepo.IsNameTaken(ctx, p.Name, dir.ParentDir); taken {
		err := errors.New(fmt.Sprint("Name already taken in destination dir. name: ", p.Name))
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	err = u.dirRepo.ChgName(ctx, uuid, p.Name)
	if err != nil {
		u.log.Err("IN [Update] failed to change name ->", err)
		err = errors.New(fmt.Sprint("Dir patch failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) Delete(c context.Context, uuid string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if constants.RootDirUuid == uuid {
		err := errors.New("Is not possible to delete the root directory")
		rErr = domain.NewUCaseErr(http.StatusConflict, err)
		return
	}

	if exists := u.dirRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	nChild, err := u.dirRepo.GetNChild(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Delete] failed to delete dir {", uuid, "} ->", err)
		err = errors.New("Failed to check directory children")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	if nChild > 0 {
		err = errors.New("Directory must be empty")
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	// TODO: Add nFiles != 0 constraint
	// TODO: Add nPlans != 0 constraint

	err = u.dirRepo.Delete(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Delete] failed to delete dir {", uuid, "} ->", err)
		err = errors.New("Failed to delete directory")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) Move(c context.Context, uuid string, nPUuid string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.dirRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if exists := u.dirRepo.ExistsByUuid(ctx, nPUuid); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if conflict := u.dirRepo.IsSubDir(ctx, uuid, nPUuid); conflict {
		err := errors.New("Reference cycle detected. Can't move dir to one of its subdirs")
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return

	}

	err := u.dirRepo.ChgParentDir(ctx, uuid, nPUuid)
	if err != nil {
		u.log.Err("IN [Move] failed to change parent dir ->", err)
		err = errors.New(fmt.Sprint("Dir patch failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *dirUsecase) Duplicate(c context.Context, p dtos.DirDuplicateDto) (res dtos.DirGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.dirRepo.ExistsByUuid(ctx, p.Uuid); !exists {
		err := errors.New("Dir not found")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if exists := u.dirRepo.ExistsByUuid(ctx, p.ParentDir); !exists {
		err := errors.New("Parent Dir not found")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if taken := u.dirRepo.IsNameTaken(ctx, p.Name, p.ParentDir); taken {
		err := errors.New(fmt.Sprint("Name taken in dest dir. name:", p.Name))
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	dirs, err := u.dirRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [Duplicate] failed to get tree ->", err)
		err = errors.New(fmt.Sprint("Dir duplication failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	branch, err := getBranchFromTree(p.Uuid, dirs)
	if err != nil {
		u.log.Err("IN [Duplicate] failed to get branch ->", err)
		err = errors.New(fmt.Sprint("Dir duplication failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	nBranch, err := flushTreeUuids(p.Uuid, branch)
	if err != nil {
		u.log.Err("IN [Duplicate] failed to flush branch uuids ->", err)
		err = errors.New(fmt.Sprint("Dir duplication failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	if len(nBranch) == 0 {
		u.log.Err("IN [Duplicate] empty new branch ->", err)
		err = errors.New(fmt.Sprint("Dir duplication failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	nBranch[0].ParentDir = p.ParentDir
	nBranch[0].Name = p.Name

	err = u.dirRepo.InsertDirs(ctx, nBranch)
	if err != nil {
		u.log.Err("IN [Duplicate] failed insert dirs ->", err)
		err = errors.New(fmt.Sprint("Dir duplication failed", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res, err = u.GetByUuid(ctx, nBranch[0].Uuid)
	if err != nil {
		u.log.Err("IN [Duplicate] failed fetch new root dir ->", err)
		err = errors.New(fmt.Sprint("Dir duplication failed", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}
