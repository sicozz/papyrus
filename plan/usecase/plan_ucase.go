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

type planUsecase struct {
	planRepo       domain.PlanRepository
	userRepo       domain.UserRepository
	dirRepo        domain.DirRepository
	taskRepo       domain.TaskRepository
	contextTimeout time.Duration
	log            utils.AggregatedLogger
}

// NewPlanUsecase will create a new pfileUsecase object representation of domain.PlanUsecase interface
func NewPlanUsecase(pr domain.PlanRepository, ur domain.UserRepository, dr domain.DirRepository, tr domain.TaskRepository, timeout time.Duration) domain.PlanUsecase {
	logger := utils.NewAggregatedLogger(constants.Usecase, constants.Plan)
	return &planUsecase{
		planRepo:       pr,
		userRepo:       ur,
		dirRepo:        dr,
		taskRepo:       tr,
		contextTimeout: timeout,
		log:            logger,
	}
}

func (u *planUsecase) GetAll(c context.Context) (res []dtos.PlanGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	plans, err := u.planRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get plans ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	plansDtos := make([]dtos.PlanGetDto, len(plans), len(plans))
	for i, p := range plans {
		plansDtos[i] = mapper.MapPlanToPlanGetDto(p)
	}

	res = plansDtos

	return
}

func (u *planUsecase) GetByUuid(c context.Context, uuid string) (res dtos.PlanGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.planRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New("Plan not found. uuid: " + uuid)
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	plan, err := u.planRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [GetByUuid] failed to get plan ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = mapper.MapPlanToPlanGetDto(plan)

	return
}

func (u *planUsecase) Store(c context.Context, dto dtos.PlanStoreDto) (res dtos.PlanGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	dateCreation, err := time.Parse(constants.LayoutDate, dto.DateCreation)
	if err != nil {
		u.log.Err("IN [Store] failed to parse DateCreation ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}
	// dateClose, err := time.Parse(constants.LayoutDate, dto.DateClose)
	// if err != nil {
	// 	u.log.Err("IN [Store] failed to parse DateClose ->", err)
	// 	rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
	// 	return
	// }

	if exists := u.dirRepo.ExistsByUuid(ctx, dto.Dir); !exists {
		err := errors.New(fmt.Sprint("Plan dir not found. uuid: ", dto.Dir))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if exists := u.userRepo.ExistsByUuid(ctx, dto.CreatorUser); !exists {
		err := errors.New(fmt.Sprint("Creator user not found. uuid: ", dto.CreatorUser))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if exists := u.userRepo.ExistsByUuid(ctx, dto.RespUser); !exists {
		err := errors.New(fmt.Sprint("Responsible user not found. uuid: ", dto.RespUser))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	p := domain.Plan{
		Code:         dto.Code,
		Name:         dto.Name,
		Origin:       dto.Origin,
		ActionType:   dto.ActionType,
		Term:         dto.Term,
		CreatorUser:  dto.CreatorUser,
		RespUser:     dto.RespUser,
		DateCreation: dateCreation,
		// DateClose: dateClose,
		DateClose:   dto.DateClose,
		Causes:      dto.Causes,
		Conclusions: dto.Conclusions,
		Dir:         dto.Dir,
		State:       dto.State,
		Stage:       dto.Stage,
		Action0desc: dto.Action0desc,
		Action0date: dto.Action0date,
		Action0user: dto.Action0user,
		Action1desc: dto.Action1desc,
		Action1date: dto.Action1date,
		Action1user: dto.Action1user,
		Action2desc: dto.Action2desc,
		Action2date: dto.Action2date,
		Action2user: dto.Action2user,
		Action3desc: dto.Action3desc,
		Action3date: dto.Action3date,
		Action3user: dto.Action3user,
		Action4desc: dto.Action4desc,
		Action4date: dto.Action4date,
		Action4user: dto.Action4user,
		Action5desc: dto.Action5desc,
		Action5date: dto.Action5date,
		Action5user: dto.Action5user,
	}

	nUuid, err := u.planRepo.Store(ctx, p)
	if err != nil {
		u.log.Err("IN [Store] failed to store plan ->", err)
		err := errors.New("Failed to store plan")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	nPlan, err := u.planRepo.GetByUuid(ctx, nUuid)
	if err != nil {
		u.log.Err("IN [Store] failed to retrieve new plan ->", err)
		err := errors.New("Failed to retrieve new plan")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	res = mapper.MapPlanToPlanGetDto(nPlan)

	respUser, err := u.userRepo.GetByUuid(ctx, nPlan.RespUser)
	if err != nil {
		u.log.Err("IN [Store] failed to retrieve responsible user ->", err)
		u.log.Err("IN [Store] failed to send task email ->", err)
		return
	}
	dirPath, err := u.dirRepo.GetPath(ctx, nPlan.Dir)
	if err != nil {
		u.log.Err("IN [Store] failed to retrieve dir ->", err)
		u.log.Err("IN [Store] failed to send task email ->", err)
		return
	}
	msg := fmt.Sprintf(
		"Se ha creado el nuevo plan de acción %v en la carpeta %v y usted ha sido designado como Responsable",
		nPlan.Name,
		dirPath,
	)
	err = utils.SendMail(respUser.Email, msg)
	if err != nil {
		u.log.Err("IN [Store] failed to send email to receiver user", err)
	}

	return
}

func (u *planUsecase) Update(c context.Context, uuid string, dto dtos.PlanUpdateDto) (res dtos.PlanGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.planRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("Plan not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	// dateClose, err := time.Parse(constants.LayoutDate, dto.DateClose)
	// if err != nil {
	// 	u.log.Err("IN [Store] failed to parse DateClose ->", err)
	// 	rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
	// 	return
	// }

	if exists := u.dirRepo.ExistsByUuid(ctx, dto.Dir); !exists {
		err := errors.New(fmt.Sprint("Plan dir not found. uuid: ", dto.Dir))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if exists := u.userRepo.ExistsByUuid(ctx, dto.RespUser); !exists {
		err := errors.New(fmt.Sprint("Responsible user not found. uuid: ", dto.RespUser))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	p := domain.Plan{
		Name:       dto.Name,
		Origin:     dto.Origin,
		ActionType: dto.ActionType,
		Term:       dto.Term,
		RespUser:   dto.RespUser,
		// DateClose:   dateClose,
		DateClose:   dto.DateClose,
		Causes:      dto.Causes,
		Conclusions: dto.Conclusions,
		Dir:         dto.Dir,
		State:       dto.State,
		Stage:       dto.Stage,
		Action0desc: dto.Action0desc,
		Action0date: dto.Action0date,
		Action0user: dto.Action0user,
		Action1desc: dto.Action1desc,
		Action1date: dto.Action1date,
		Action1user: dto.Action1user,
		Action2desc: dto.Action2desc,
		Action2date: dto.Action2date,
		Action2user: dto.Action2user,
		Action3desc: dto.Action3desc,
		Action3date: dto.Action3date,
		Action3user: dto.Action3user,
		Action4desc: dto.Action4desc,
		Action4date: dto.Action4date,
		Action4user: dto.Action4user,
		Action5desc: dto.Action5desc,
		Action5date: dto.Action5date,
		Action5user: dto.Action5user,
	}

	err := u.planRepo.Update(ctx, uuid, p)
	if err != nil {
		u.log.Err("IN [Update] failed to update plan ->", err)
		err := errors.New("Failed to update plan")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	nPlan, err := u.planRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Update] failed to retrieve plan ->", err)
		err := errors.New("Failed to retrieve plan")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	res = mapper.MapPlanToPlanGetDto(nPlan)

	if dto.State != "revisado" {
		return
	}

	respUser, err := u.userRepo.GetByUuid(ctx, nPlan.RespUser)
	if err != nil {
		u.log.Err("IN [Store] failed to retrieve responsible user ->", err)
		u.log.Err("IN [Store] failed to send task email ->", err)
		return
	}
	dirPath, err := u.dirRepo.GetPath(ctx, nPlan.Dir)
	if err != nil {
		u.log.Err("IN [Store] failed to retrieve dir ->", err)
		u.log.Err("IN [Store] failed to send task email ->", err)
		return
	}
	msg := fmt.Sprintf(
		"El plan %v en la carpeta %v se revisó y está esperando para ser Cerrado",
		nPlan.Name,
		dirPath,
	)
	err = utils.SendMail(respUser.Email, msg)
	if err != nil {
		u.log.Err("IN [Store] failed to send email to receiver user", err)
	}

	return
}

func (u *planUsecase) Delete(c context.Context, uuid string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.planRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("Plan not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	err := u.planRepo.Delete(ctx, uuid)
	if err != nil {
		u.log.Err(fmt.Sprintf("IN [Delete] failed to delete plan uuid:%v -> %v", uuid, err))
		err = errors.New("Failed to delete plan")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}
