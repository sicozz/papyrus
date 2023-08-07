package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sicozz/papyrus/domain"
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
	// Refactor filluserdetails
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistByUname(ctx, uname); !exists {
		err = errors.New(fmt.Sprint("User not found. username: ", uname))
		return
	}

	res, err = u.userRepo.GetByUsername(ctx, uname)
	if err != nil {
		domain.AgLog.Error("Error inside Fetch function", err)
		err = errors.New(fmt.Sprint("User fetch failed. username: ", uname))
		return domain.User{}, err
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

	if exists := u.userRepo.ExistByUname(ctx, user.Username); exists {
		err = errors.New("Username already taken")
		return
	}

	if exists := u.userRepo.ExistByEmail(ctx, user.Email); exists {
		err = errors.New("Email already taken")
		return
	}

	r, err := u.roleRepo.GetByDescription(ctx, defRoleDesc)
	if err != nil {
		domain.AgLog.Error("Could not find default role.", err)
		err = errors.New("Base role fetch failed")
		return
	}
	user.Role = r

	s, err := u.userStateRepo.GetByDescription(ctx, defUserStateDesc)
	if err != nil {
		domain.AgLog.Error("Base state fetch failed", err)
		err = errors.New("Base role fetch failed")
		return
	}
	user.State = s

	err = u.userRepo.Store(ctx, user)
	return
}

func (u *userUsecase) Delete(c context.Context, uname string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistByUname(ctx, uname); !exists {
		err = errors.New(fmt.Sprint("User not found. username: ", uname))
		return
	}

	err = u.userRepo.Delete(ctx, uname)
	if err != nil {
		domain.AgLog.Error("User deletion failed. username: ", uname)
	}

	return
}

func (u *userUsecase) Update(c context.Context, uname string, uUp *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if uUp.Email != "" {
		if taken := u.userRepo.ExistByEmail(ctx, uUp.Email); taken {
			err = errors.New("Email already taken")
			return
		}
		err = u.userRepo.ChgEmail(ctx, uname, uUp.Email)
		if err != nil {
			return errors.New(fmt.Sprint("User patch failed: ", err))
		}
	}

	if uUp.Name != "" {
		err = u.userRepo.ChgName(ctx, uname, uUp.Name)
		if err != nil {
			return errors.New(fmt.Sprint("User patch failed: ", err))
		}
	}

	if uUp.Lastname != "" {
		err = u.userRepo.ChgLstname(ctx, uname, uUp.Lastname)
		if err != nil {
			return errors.New(fmt.Sprint("User patch failed: ", err))
		}
	}

	if uUp.Role.Description != "" {
		r, rErr := u.roleRepo.GetByDescription(ctx, uUp.Role.Description)
		if rErr != nil {
			domain.AgLog.Error("Role not found: ", rErr)
			return errors.New(fmt.Sprint("Role not found"))
		}

		err = u.userRepo.ChgRole(ctx, uname, r)
		if err != nil {
			return errors.New(fmt.Sprint("User patch failed: ", err))
		}
	}

	if uUp.State.Description != "" {
		s, sErr := u.userStateRepo.GetByDescription(ctx, uUp.State.Description)
		if sErr != nil {
			domain.AgLog.Error("User_state not found", sErr)
			return errors.New(fmt.Sprint("User_state not found"))
		}

		err = u.userRepo.ChgState(ctx, uname, s)
		if err != nil {
			return errors.New(fmt.Sprint("User patch failed: ", err))
		}
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
