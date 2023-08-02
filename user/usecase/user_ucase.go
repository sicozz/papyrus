package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
)

const (
	defRoleDesc      = `estandar`
	defUserStateDesc = `inactivo`
)

type userUsecase struct {
	userRepo       domain.UserRepository
	roleRepo       domain.RoleRepository
	userStateRepo  domain.UserStateRepository
	contextTimeout time.Duration
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (u *userUsecase) fillUserDetails(ctx context.Context, users []domain.User) (err error) {
	// get roles
	roles, err := u.roleRepo.GetAll(ctx)
	if err != nil {
		domain.AgLog.Error("Could not get roles to fill user details", err)
	}

	mapRoles := map[int64]domain.Role{}
	for _, role := range roles { //nolint
		mapRoles[role.Code] = role
	}

	// get user_states
	states, err := u.userStateRepo.GetAll(ctx)
	if err != nil {
		domain.AgLog.Error("Could not get user_states to fill user details", err)
	}

	mapStates := map[int64]domain.UserState{}
	for _, state := range states { //nolint
		mapStates[state.Code] = state
	}

	// merge the user's data
	for idx, user := range users { //nolint
		if r, ok := mapRoles[user.Role.Code]; ok {
			users[idx].Role = r
		}

		if s, ok := mapStates[user.State.Code]; ok {
			users[idx].State = s
		}
	}

	return
}

// NewUserUsecase will create a new userUsecase object representation of domain.UserUsecase interface
func NewUserUsecase(ur domain.UserRepository, rr domain.RoleRepository, usr domain.UserStateRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       ur,
		roleRepo:       rr,
		userStateRepo:  usr,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Fetch(c context.Context) (res []domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.userRepo.Fetch(ctx)
	if err != nil {
		domain.AgLog.Error("Error inside Fetch function")
	}

	err = u.fillUserDetails(ctx, res)
	if err != nil {
		domain.AgLog.Error("Error filling user details")
	}
	return
}

func (u *userUsecase) GetByUsername(c context.Context, uname string) (res domain.User, err error) {
	// Refactor flluserdetails
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.userRepo.GetByUsername(ctx, uname)
	if err != nil {
		domain.AgLog.Error("Error inside Fetch function")
	}

	resArr := make([]domain.User, 1)
	resArr[0] = res
	err = u.fillUserDetails(ctx, resArr)
	res = resArr[0]
	if err != nil {
		domain.AgLog.Error("Error filling user details")
	}
	return
}

func (u *userUsecase) Store(c context.Context, user *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	r, err := u.roleRepo.GetByDescription(ctx, defRoleDesc)
	if err != nil {
		domain.AgLog.Error("Could not find default role.", err)
	}
	user.Role = r

	s, err := u.userStateRepo.GetByDescription(ctx, defUserStateDesc)
	if err != nil {
		domain.AgLog.Error("Could not find default user state.", err)
	}
	user.State = s

	err = u.userRepo.Store(ctx, user)
	return
}

func (u *userUsecase) Delete(c context.Context, uname string) (body dtos.BaseDto, errBody dtos.ErrorDto) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err := u.userRepo.Delete(ctx, uname)
	if err != nil {
		domain.AgLog.Error("error while deleting user with username:", uname, err)
		errBody = dtos.NewErrorDto(fmt.Sprint("fail: user deletion [username:", uname, "]"))
		return
	}

	body = dtos.NewBaseDto(fmt.Sprint("success: User deletion [username:", uname, "]"))
	return
}

func (u *userUsecase) ChangeState(c context.Context, uname string, desc string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	s, err := u.userStateRepo.GetByDescription(ctx, desc)
	if err != nil {
		domain.AgLog.Error("Could not find user state with description:", desc)
		return
	}

	err = u.userRepo.ChangeState(ctx, uname, s)
	if err != nil {
		domain.AgLog.Error("Could not update state of user:", uname)
	}

	return
}

func (u *userUsecase) ChangeRole(c context.Context, uname string, desc string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	r, err := u.roleRepo.GetByDescription(ctx, desc)
	if err != nil {
		domain.AgLog.Error("Could not find role with description:", desc)
		return
	}

	err = u.userRepo.ChangeRole(ctx, uname, r)
	if err != nil {
		domain.AgLog.Error("Could not update role of user:", uname)
	}

	return
}

func (u *userUsecase) Login(c context.Context, uname string, passwd string) (res domain.User, err error) {
	// Refactor flluserdetails
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.userRepo.Login(ctx, uname, passwd)
	if err != nil {
		domain.AgLog.Error("Error inside Login function")
		return
	}

	resArr := make([]domain.User, 1)
	resArr[0] = res
	err = u.fillUserDetails(ctx, resArr)
	res = resArr[0]
	if err != nil {
		domain.AgLog.Error("Error filling user details")
	}
	return
}
