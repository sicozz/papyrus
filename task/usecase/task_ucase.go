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

type taskUsecase struct {
	taskRepo       domain.TaskRepository
	dirRepo        domain.DirRepository
	userRepo       domain.UserRepository
	contextTimeout time.Duration
	log            utils.AggregatedLogger
}

// NewTaskUsecase will create a new taskUsecase object representation of domain.TaskUsecase interface
func NewTaskUsecase(tr domain.TaskRepository, dr domain.DirRepository, ur domain.UserRepository, timeout time.Duration) domain.TaskUsecase {
	logger := utils.NewAggregatedLogger(constants.Usecase, constants.Task)
	return &taskUsecase{
		taskRepo:       tr,
		dirRepo:        dr,
		userRepo:       ur,
		contextTimeout: timeout,
		log:            logger,
	}
}

func (u *taskUsecase) GetAll(c context.Context) (res []dtos.TaskGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	tasks, err := u.taskRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get tasks ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	tasksDtos := make([]dtos.TaskGetDto, len(tasks), len(tasks))
	for i, t := range tasks {
		tasksDtos[i] = mapper.MapTaskToTaskGetDto(t)
	}

	res = tasksDtos

	return
}

func (u *taskUsecase) GetByUuid(c context.Context, uuid string) (res dtos.TaskGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.taskRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New("Task not found. uuid: " + uuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	task, err := u.taskRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetByUuid] failed to get task ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = mapper.MapTaskToTaskGetDto(task)

	return
}

func (u *taskUsecase) GetByUser(c context.Context, uuid string) (res []dtos.TaskGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New("User not found. uuid: " + uuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	tasks, err := u.taskRepo.GetByUser(ctx, uuid)
	if err != nil {
		u.log.Err(fmt.Sprintf("IN [GetByUser] failed to get tasks from user uuid: %v -> %v", uuid, err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = make([]dtos.TaskGetDto, len(tasks), len(tasks))
	for i, t := range tasks {
		res[i] = mapper.MapTaskToTaskGetDto(t)
	}

	return
}

func (u *taskUsecase) GetByPlan(c context.Context, uuid string) (res []dtos.TaskGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	tasks, err := u.taskRepo.GetByPlan(ctx, uuid)
	if err != nil {
		u.log.Err(fmt.Sprintf("IN [GetByPlan] failed to get tasks from plan uuid: %v -> %v", uuid, err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = make([]dtos.TaskGetDto, len(tasks), len(tasks))
	for i, t := range tasks {
		res[i] = mapper.MapTaskToTaskGetDto(t)
	}

	return
}

func (u *taskUsecase) Store(c context.Context, p dtos.TaskStoreDto) (res dtos.TaskGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// parse date_creation
	dateCreation, err := time.Parse(constants.LayoutDate, p.DateCreate)
	if err != nil {
		u.log.Err("IN [Store] failed to parse DateCreation ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	// check dir
	if exists := u.dirRepo.ExistsByUuid(ctx, p.Dir); !exists {
		err := errors.New(fmt.Sprint("Dir not found. uuid: ", p.Dir))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	// check creator user
	if exists := u.userRepo.ExistsByUuid(ctx, p.CreatorUser); !exists {
		err := errors.New(fmt.Sprint("Creator user not found. uuid: ", p.CreatorUser))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	// check receiver user
	if exists := u.userRepo.ExistsByUuid(ctx, p.RecvUser); !exists {
		err := errors.New(fmt.Sprint("Receiver user not found. uuid: ", p.RecvUser))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	t := domain.Task{
		Name:         p.Name,
		Procedure:    p.Procedure,
		DateCreation: dateCreation,
		Term:         p.Term,
		Dir:          p.Dir,
		CreatorUser:  p.CreatorUser,
		RecvUser:     p.RecvUser,
		Plan:         p.Plan,
	}

	nUuid, err := u.taskRepo.Store(ctx, t)
	if err != nil {
		u.log.Err("IN [Store] failed to store task ->", err)
		err := errors.New("Failed to store task")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	nTask, err := u.taskRepo.GetByUuid(ctx, nUuid)
	if err != nil {
		u.log.Err("IN [Store] failed to retrieve new task ->", err)
		err := errors.New("Failed to retrieve new task")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	res = mapper.MapTaskToTaskGetDto(nTask)

	creatorUser, err := u.userRepo.GetByUuid(ctx, p.CreatorUser)
	if err != nil {
		u.log.Err("IN [Store] failed to retrieve creator user ->", err)
		u.log.Err("IN [Store] failed to send task email ->", err)
		return
	}
	rcvUser, err := u.userRepo.GetByUuid(ctx, p.RecvUser)
	if err != nil {
		u.log.Err("IN [Store] failed to retrieve responsible user ->", err)
		u.log.Err("IN [Store] failed to send task email ->", err)
		return
	}
	dirPath, err := u.dirRepo.GetPath(ctx, nTask.Dir)
	if err != nil {
		u.log.Err("IN [Store] failed to retrieve dir ->", err)
		u.log.Err("IN [Store] failed to send task email ->", err)
		return
	}
	msg := fmt.Sprintf(
		"Este es un mensaje del sistema Papyrus.\nEl usuario %v le ha asignado el dia %v la tarea %v en la carpeta %v, con las siguientes actividades:\n%v",
		creatorUser.Username,
		nTask.DateCreation.Format("Monday, January 2, 2006 15:04:05"),
		nTask.Name,
		dirPath,
		nTask.Procedure,
	)
	u.log.Inf("Mail msg:", msg)
	go utils.SendMail(rcvUser.Email, msg)

	// err = utils.SendMail(rcvUser.Email, msg)
	// if err != nil {
	// 	u.log.Err("IN [Store] failed to send email to receiver user", err)
	// }

	return
}

func (u *taskUsecase) StoreMultiple(c context.Context, inputDtos []dtos.TaskStoreDto) (res []dtos.TaskGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	tasks := []domain.Task{}
	for _, iDto := range inputDtos {
		// parse date_creation
		dateCreation, err := time.Parse(constants.LayoutDate, iDto.DateCreate)
		if err != nil {
			u.log.Err("IN [StoreMultiple] failed to parse DateCreation ->", err)
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}

		// check dir
		if exists := u.dirRepo.ExistsByUuid(ctx, iDto.Dir); !exists {
			err := errors.New(fmt.Sprint("Dir not found. uuid: ", iDto.Dir))
			rErr = domain.NewUCaseErr(http.StatusNotFound, err)
			return
		}

		// check creator user
		if exists := u.userRepo.ExistsByUuid(ctx, iDto.CreatorUser); !exists {
			err := errors.New(fmt.Sprint("Creator user not found. uuid: ", iDto.CreatorUser))
			rErr = domain.NewUCaseErr(http.StatusNotFound, err)
			return
		}

		// check receiver user
		if exists := u.userRepo.ExistsByUuid(ctx, iDto.RecvUser); !exists {
			err := errors.New(fmt.Sprint("Receiver user not found. uuid: ", iDto.RecvUser))
			rErr = domain.NewUCaseErr(http.StatusNotFound, err)
			return
		}

		t := domain.Task{
			Name:         iDto.Name,
			Procedure:    iDto.Procedure,
			DateCreation: dateCreation,
			Term:         iDto.Term,
			Dir:          iDto.Dir,
			CreatorUser:  iDto.CreatorUser,
			RecvUser:     iDto.RecvUser,
			Plan:         iDto.Plan,
		}

		tasks = append(tasks, t)
	}

	err := u.taskRepo.StoreMultiple(ctx, tasks)
	if err != nil {
		u.log.Err("IN [StoreMultiple] failed to store tasks ->", err)
		err := errors.New("Failed to store tasks")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	nTasks, err := u.taskRepo.GetByPlan(ctx, tasks[0].Plan)
	if err != nil {
		u.log.Err("IN [StoreMultiple] failed to retrieve new tasks ->", err)
		err := errors.New("Failed to retrieve new tasks")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	res = []dtos.TaskGetDto{}
	for _, nT := range nTasks {
		res = append(res, mapper.MapTaskToTaskGetDto(nT))
		rcvUser, err := u.userRepo.GetByUuid(ctx, nT.RecvUser)
		if err != nil {
			u.log.Err("IN [Store] failed to retrieve responsible user ->", err)
			u.log.Err("IN [Store] failed to send task email ->", err)
			return
		}
		dirPath, err := u.dirRepo.GetPath(ctx, nT.Dir)
		if err != nil {
			u.log.Err("IN [Store] failed to retrieve dir ->", err)
			u.log.Err("IN [Store] failed to send task email ->", err)
			return
		}
		msg := fmt.Sprintf(
			"Se le ha asignado la tarea %v en la carpeta %v",
			nT.Name,
			dirPath,
		)
		go utils.SendMail(rcvUser.Email, msg)

		// err = utils.SendMail(rcvUser.Email, msg)
		// if err != nil {
		// 	u.log.Err("IN [Store] failed to send email to receiver user", err)
		// }
	}

	return
}

func (u *taskUsecase) ChgCheck(c context.Context, tUuid, uUuid string, chk bool) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// check task
	if exists := u.taskRepo.ExistsByUuid(ctx, tUuid); !exists {
		err := errors.New(fmt.Sprint("Task not found. uuid: ", uUuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	// check user
	if exists := u.userRepo.ExistsByUuid(ctx, uUuid); !exists {
		err := errors.New(fmt.Sprint("User not found. uuid: ", uUuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	t, err := u.taskRepo.GetByUuid(ctx, tUuid)
	if err != nil {
		u.log.Err("IN [ChgState] failed to fetch task. uuid ", tUuid)
		err := errors.New(fmt.Sprint("Failed to fetch task"))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if t.RecvUser != uUuid {
		err := errors.New(fmt.Sprintf("User %v must be reciever user of task %v", uUuid, tUuid))
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	err = u.taskRepo.ChgCheck(ctx, tUuid, uUuid, chk)
	if err != nil {
		u.log.Err("IN [ChgCheck] failed to check task", tUuid, " with user ", uUuid, " -> ", err)
		err = errors.New("Failed to check task")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	if chk == true {
		err = u.taskRepo.SetDateCheck(ctx, tUuid)
		if err != nil {
			u.log.Err("IN [Store] failed to set date_check for task ->", err)
			return
		}
		creatorUser, err := u.userRepo.GetByUuid(ctx, t.CreatorUser)
		if err != nil {
			u.log.Err("IN [Store] failed to retrieve responsible user ->", err)
			u.log.Err("IN [Store] failed to send task email ->", err)
			return
		}
		dirPath, err := u.dirRepo.GetPath(ctx, t.Dir)
		if err != nil {
			u.log.Err("IN [Store] failed to retrieve dir ->", err)
			u.log.Err("IN [Store] failed to send task email ->", err)
			return
		}
		msg := fmt.Sprintf(
			"La tarea %v en la carpeta %v se ejecutó y está esperando para ser Cumplida.",
			t.Name,
			dirPath,
		)
		go utils.SendMail(creatorUser.Email, msg)

		// err = utils.SendMail(rcvUser.Email, msg)
		// if err != nil {
		// 	u.log.Err("IN [Store] failed to send email to receiver user", err)
		// }
	}

	return
}

func (u *taskUsecase) ChgState(c context.Context, tUuid, uUuid string, desc string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// check task_state
	if exists := u.taskRepo.ExistsStateByDesc(ctx, desc); !exists {
		err := errors.New(fmt.Sprint("Task state not found. description: ", desc))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	// check task
	if exists := u.taskRepo.ExistsByUuid(ctx, tUuid); !exists {
		err := errors.New(fmt.Sprint("Task not found. uuid: ", uUuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	// check user
	if exists := u.userRepo.ExistsByUuid(ctx, uUuid); !exists {
		err := errors.New(fmt.Sprint("User not found. uuid: ", uUuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	t, err := u.taskRepo.GetByUuid(ctx, tUuid)
	if err != nil {
		u.log.Err("IN [ChgState] failed to fetch task. uuid ", tUuid)
		err := errors.New(fmt.Sprint("Failed to fetch task"))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if t.CreatorUser != uUuid {
		err := errors.New(fmt.Sprintf("User %v must be creator user of task %v", uUuid, tUuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	err = u.taskRepo.ChgState(ctx, tUuid, uUuid, desc)
	if err != nil {
		u.log.Err("IN [ChgState] failed to change task state of task ", tUuid, " with user ", uUuid, " -> ", err)
		err = errors.New("Failed to change task state")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *taskUsecase) Delete(c context.Context, tUuid, uUuid string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.taskRepo.ExistsByUuid(ctx, tUuid); !exists {
		err := errors.New(fmt.Sprint("Task not found. uuid: ", tUuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	task, err := u.taskRepo.GetByUuid(ctx, tUuid)
	if err != nil {
		u.log.Err(fmt.Sprintf("IN [Delete] failed fetch task uuid:%v -> %v", tUuid, err))
		err = errors.New("Failed to delete task")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	if task.CreatorUser != uUuid {
		err = errors.New("User must be the task creator to delete the task")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	if exists := u.userRepo.ExistsByUuid(ctx, uUuid); !exists {
		err := errors.New(fmt.Sprint("User not found. uuid: ", uUuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	err = u.taskRepo.Delete(ctx, tUuid, uUuid)
	if err != nil {
		u.log.Err(fmt.Sprintf("IN [Delete] failed to delete task uuid:%v with user uuid:%v -> %v", tUuid, uUuid, err))
		err = errors.New("Failed to delete task")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}
