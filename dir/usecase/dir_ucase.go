package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/domain/mapper"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type dirUsecase struct {
	dirRepo        domain.DirRepository
	userRepo       domain.UserRepository
	pFileRepo      domain.PFileRepository
	taskRepo       domain.TaskRepository
	planRepo       domain.PlanRepository
	contextTimeout time.Duration
	log            utils.AggregatedLogger
}

// NewDirUsecase will create a new dirUsecase object representation of domain.DirUsecase interface
func NewDirUsecase(dr domain.DirRepository, ur domain.UserRepository, pfr domain.PFileRepository, tr domain.TaskRepository, pr domain.PlanRepository, timeout time.Duration) domain.DirUsecase {
	logger := utils.NewAggregatedLogger(constants.Usecase, constants.Dir)
	return &dirUsecase{
		dirRepo:        dr,
		userRepo:       ur,
		pFileRepo:      pfr,
		taskRepo:       tr,
		planRepo:       pr,
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

	plans, err := u.planRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get plan ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}
	for _, p := range plans {
		planDirDto := mapper.MapPlanToDir(p)
		dirs = append(dirs, planDirDto)
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
				// WARN: Esto es demasiado mediocre
				if pf.Subtype == "registro" || pf.Subtype == "evidencia" {
					res[j].Type = pf.Subtype
				}
				res[j].State = pf.State

				res[j].RespUser = pf.RespUser
				res[j].Subtype = pf.Subtype
				res[j].Datecreate = pf.DateCreation.Format(constants.LayoutDate)
				res[j].Term = pf.Term
			}
		}
	}
	for _, t := range tasks {
		for j := range res {
			if t.Uuid == res[j].Uuid {
				// TODO: put this type and state in mapper
				res[j].Type = "tarea"
				res[j].State = t.State
				res[j].Term = t.Term
			}
		}
	}
	for _, p := range plans {
		for j := range res {
			if p.Uuid == res[j].Uuid {
				// TODO: put this type and state in mapper
				res[j].Type = "plan"
				res[j].State = p.State
				res[j].CreatorUser = p.CreatorUser
				res[j].RespUser = p.RespUser
				res[j].Datecreate = p.DateCreation.Format(constants.LayoutDate)
				res[j].DateClose = p.DateClose
				res[j].Term = p.Term
			}
		}
	}

	return
}

func (u *dirUsecase) GetByUuid(c context.Context, uuid string) (res dtos.DirGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.dirRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New("Dir not found")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

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

func (u *dirUsecase) GetDocsByUser(c context.Context, uuid string) (res []dtos.DocsNotDirGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New("User not found. uuid: " + uuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	tasks, err := u.taskRepo.GetByUser(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetDocsByUser] failed to get tasks ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	pFiles, err := u.pFileRepo.GetByUser(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetDocsByUser] failed to get pfiles ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	plans, err := u.planRepo.GetByUser(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetDocsByUser] failed to get plans ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = []dtos.DocsNotDirGetDto{}
	for _, pf := range pFiles {
		apps := []domain.Approvation{}
		if pf.Subtype != "registro" {
			apps, err = u.pFileRepo.GetApprovations(ctx, pf.Uuid)
			if err != nil {
				u.log.Err("IN [GetDocsByUser] failed to get file approvations ->", err)
				rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
				return
			}
		}

		dnd := mapper.MapPFileToDocsNotDirGetDto(pf, apps)
		dnd.Type = "documento"
		dnd.Path, err = u.dirRepo.GetPath(ctx, dnd.ParentDir)
		if err != nil {
			u.log.Err("IN [GetDocsByUser] failed to get file path ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
		dnd.Depth, err = u.dirRepo.GetDepth(ctx, dnd.ParentDir)
		if err != nil {
			u.log.Err("IN [GetDocsByUser] failed to get file depth ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}

		res = append(res, dnd)
	}

	for _, t := range tasks {
		dnd := mapper.MapTaskToDocsNotDirGetDto(t)
		dnd.Type = "tarea"
		dnd.Path, err = u.dirRepo.GetPath(ctx, dnd.ParentDir)
		if err != nil {
			u.log.Err("IN [GetDocsByUser] failed to get task path ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
		dnd.Depth, err = u.dirRepo.GetDepth(ctx, dnd.ParentDir)
		if err != nil {
			u.log.Err("IN [GetDocsByUser] failed to get task depth ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
		res = append(res, dnd)
	}

	for _, p := range plans {
		dnd := mapper.MapPlanToDocsNotDirGetDto(p)
		dnd.Type = "plan"
		dnd.Path, err = u.dirRepo.GetPath(ctx, dnd.ParentDir)
		if err != nil {
			u.log.Err("IN [GetDocsByUser] failed to get plan path ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
		dnd.Depth, err = u.dirRepo.GetDepth(ctx, dnd.ParentDir)
		if err != nil {
			u.log.Err("IN [GetDocsByUser] failed to get plan depth ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
		res = append(res, dnd)
	}

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

	dirs, err := u.dirRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [Delete] failed to get dirs ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	path, err := u.dirRepo.GetPath(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Delete] failed to get dir path ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	nChild, err := u.dirRepo.GetNChild(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Delete] failed to get dir children number ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	depth, err := u.dirRepo.GetDepth(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Delete] failed to get dir depth ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	subtree, err := fillDetailsBFS(uuid, dirs, path, nChild, depth)
	if err != nil {
		u.log.Err("IN [Delete] failed to fill tree details ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	for i := len(subtree) - 1; i >= 0; i-- {
		err = u.dirRepo.Delete(ctx, subtree[i].Uuid)
		if err != nil {
			u.log.Err("IN [Delete] failed to delete dir {", uuid, "} ->", err)
			err = errors.New("Failed to delete directory")
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
	}

	// No checks for dir deletion since client knows better (for a qas...)
	// nChild, err := u.dirRepo.GetNChild(ctx, uuid)
	// if err != nil {
	// 	u.log.Err("IN [Delete] failed to delete dir {", uuid, "} ->", err)
	// 	err = errors.New("Failed to check directory children")
	// 	rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
	// 	return
	// }

	// if nChild > 0 {
	// 	err = errors.New("Directory must be empty")
	// 	rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
	// 	return
	// }

	// err = u.dirRepo.Delete(ctx, uuid)
	// if err != nil {
	// 	u.log.Err("IN [Delete] failed to delete dir {", uuid, "} ->", err)
	// 	err = errors.New("Failed to delete directory")
	// 	rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
	// 	return
	// }

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

func (u *dirUsecase) Duplicate(c context.Context, p dtos.DirDuplicateDto) (res []dtos.DirGetDto, rErr domain.RequestErr) {
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

	nBranchDtos := []dtos.DirGetDto{}
	for _, d := range nBranch {
		path, err := u.dirRepo.GetPath(ctx, d.Uuid)
		if err != nil {
			u.log.Err("IN [Duplicate] failed to get dir path ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}

		nChild, err := u.dirRepo.GetNChild(ctx, d.Uuid)
		if err != nil {
			u.log.Err("IN [Duplicate] failed to get dir children number ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}

		depth, err := u.dirRepo.GetDepth(ctx, d.Uuid)
		if err != nil {
			u.log.Err("IN [Duplicate] failed to get dir depth ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}

		dto := mapper.MapDirToDirGetDto(d, path, nChild, depth)

		nBranchDtos = append(nBranchDtos, dto)
	}

	res = nBranchDtos

	return
}

func (u *dirUsecase) GetDirSize(c context.Context, uuid string) (res dtos.DirSizeGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.dirRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New("Dir not found")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	dirs, err := u.dirRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetDirSize] failed to get dirs ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	path, err := u.dirRepo.GetPath(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetDirSize] failed to get dir path ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	nChild, err := u.dirRepo.GetNChild(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetDirSize] failed to get dir children number ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	depth, err := u.dirRepo.GetDepth(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetDirSize] failed to get dir depth ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	subtree, err := fillDetailsBFS(uuid, dirs, path, nChild, depth)
	if err != nil {
		u.log.Err("IN [GetDirSize] failed to fill tree details ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	dirPfiles := []domain.PFile{}
	for _, d := range subtree {
		subdirPFiles, err := u.pFileRepo.GetByDir(ctx, d.Uuid)
		if err != nil {
			u.log.Err("IN [GetDirSize] failed to get subdir files ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
		dirPfiles = append(dirPfiles, subdirPFiles...)
	}

	var dirSize int64 = 0
	for _, dPF := range dirPfiles {
		fileInfo, err := os.Stat(dPF.FsPath)
		if err != nil {
			u.log.Err("IN [GetDirSize] failed to get file size ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
		dirSize += fileInfo.Size()
		res.FileCount += 1
	}

	res.Size = fmt.Sprintf("%.2f GB", float64(dirSize)/(1024*1024*1024))

	return
}

func (u *dirUsecase) AddRecursivePermission(c context.Context, d dtos.UserAddPermissionDto) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistsByUuid(ctx, d.UserUuid); !exists {
		err := errors.New("User not found. uuid:" + d.UserUuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if exists := u.dirRepo.ExistsByUuid(ctx, d.DirUuid); !exists {
		err := errors.New("Dir not found. uuid:" + d.DirUuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	dirs, err := u.dirRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [Delete] failed to get dirs ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	path, err := u.dirRepo.GetPath(ctx, d.DirUuid)
	if err != nil {
		u.log.Err("IN [Delete] failed to get dir path ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	nChild, err := u.dirRepo.GetNChild(ctx, d.DirUuid)
	if err != nil {
		u.log.Err("IN [Delete] failed to get dir children number ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	depth, err := u.dirRepo.GetDepth(ctx, d.DirUuid)
	if err != nil {
		u.log.Err("IN [Delete] failed to get dir depth ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	subtree, err := fillDetailsBFS(d.DirUuid, dirs, path, nChild, depth)
	if err != nil {
		u.log.Err("IN [Delete] failed to fill tree details ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	for _, subdir := range subtree {
		if exists := u.userRepo.ExistsPermission(ctx, d.UserUuid, subdir.Uuid); exists {
			continue
		}

		p := domain.Permission{UserUuid: d.UserUuid, DirUuid: subdir.Uuid}

		err := u.userRepo.AddPermission(ctx, p)
		if err != nil {
			u.log.Err("IN [AddPermission] failed to add permission -> ", err)
			rErr = domain.NewUCaseErr(http.StatusNotFound, err)
			return
		}
	}

	return
}
