package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
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
	log            utils.AggregatedLogger
}

// NewUserUsecase will create a new userUsecase object representation of domain.UserUsecase interface
func NewUserUsecase(ur domain.UserRepository, rr domain.RoleRepository, usr domain.UserStateRepository, timeout time.Duration) domain.UserUsecase {
	logger := utils.NewAggregatedLogger(constants.Usecase, constants.User)
	return &userUsecase{
		userRepo:       ur,
		roleRepo:       rr,
		userStateRepo:  usr,
		contextTimeout: timeout,
		log:            logger,
	}
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
		u.log.Err("IN [fillUserDetails] failed to get roles ->", err)
	}

	mapRoles := map[int64]domain.Role{}
	for _, role := range roles { //nolint
		mapRoles[role.Code] = role
	}

	// get user_states
	states, err := u.userStateRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [fillUserDetails] failed to get user_states ->", err)
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

func (u *userUsecase) GetAll(c context.Context) (res []domain.User, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get users ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = users
	return
}

func (u *userUsecase) GetByUsername(c context.Context, uname string) (res domain.User, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistByUname(ctx, uname); !exists {
		err := errors.New(fmt.Sprint("User not found. username: ", uname))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	res, err := u.userRepo.GetByUsername(ctx, uname)
	if err != nil {
		u.log.Err("IN [GetByUsername] failed to get user ->", err)
		err = errors.New(fmt.Sprint("User fetch failed. username: ", uname))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return domain.User{}, rErr
	}

	return
}

func (u *userUsecase) Store(c context.Context, user *domain.User) (res domain.User, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistByUname(ctx, user.Username); exists {
		err := errors.New("Username already taken")
		rErr = domain.NewUCaseErr(http.StatusConflict, err)
		return
	}

	if exists := u.userRepo.ExistByEmail(ctx, user.Email); exists {
		err := errors.New("Email already taken")
		rErr = domain.NewUCaseErr(http.StatusConflict, err)
		return
	}

	r, err := u.roleRepo.GetByDescription(ctx, defRoleDesc)
	if err != nil {
		err = errors.New("Base role fetch failed")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}
	user.Role = r

	s, err := u.userStateRepo.GetByDescription(ctx, defUserStateDesc)
	if err != nil {
		err = errors.New("Base role fetch failed")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}
	user.State = s

	err = u.userRepo.Store(ctx, user)
	if err != nil {
		err = errors.New("User creation failed")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	nUser, err := u.userRepo.GetByUsername(ctx, user.Username)
	if err != nil {
		err = errors.New("User creation failed")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = nUser

	return
}

func (u *userUsecase) Delete(c context.Context, uname string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistByUname(ctx, uname); !exists {
		err := errors.New(fmt.Sprint("User not found. username: ", uname))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	err := u.userRepo.Delete(ctx, uname)
	if err != nil {
		u.log.Err("IN [Delete] failed to delete user {", uname, "} ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *userUsecase) Update(c context.Context, uuid string, uUp dtos.UserUpdateDto) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("User not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	user, err := u.userRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Update] failed to get user ->", err)
		err = errors.New(fmt.Sprint("User patch failed. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}
	uname := user.Username

	if uUp.Email != "" {
		err := u.userRepo.ChgEmail(ctx, uname, uUp.Email)
		if err != nil {
			u.log.Err("IN [Update] failed to change email ->", err)
			err = errors.New(fmt.Sprint("User patch failed: ", err))
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
	}

	if uUp.Name != "" {
		err := u.userRepo.ChgName(ctx, uname, uUp.Name)
		if err != nil {
			u.log.Err("IN [Update] failed to change name ->", err)
			err = errors.New(fmt.Sprint("User patch failed: ", err))
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
	}

	if uUp.Lastname != "" {
		err := u.userRepo.ChgLstname(ctx, uname, uUp.Lastname)
		if err != nil {
			u.log.Err("IN [Update] failed to change lastname ->", err)
			err = errors.New(fmt.Sprint("User patch failed: ", err))
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
	}

	if uUp.Role != "" {
		r, err := u.roleRepo.GetByDescription(ctx, uUp.Role)
		if err != nil {
			u.log.Err("IN [Update] failed to get role ->", err)
			err = errors.New(fmt.Sprint("Role not found"))
			rErr = domain.NewUCaseErr(http.StatusNotFound, err)
			return
		}
		err = u.userRepo.ChgRole(ctx, uname, r)
		if err != nil {
			u.log.Err("IN [Update] failed to change role ->", err)
			err = errors.New(fmt.Sprint("User patch failed: ", err))
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
	}

	if uUp.State != "" {
		s, err := u.userStateRepo.GetByDescription(ctx, uUp.State)
		if err != nil {
			u.log.Err("IN [Update] failed to get user_state ->", err)
			err = errors.New(fmt.Sprint("User_state not found"))
			rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		}

		err = u.userRepo.ChgState(ctx, uname, s)
		if err != nil {
			u.log.Err("IN [Update] failed to change user_state ->", err)
			err = errors.New(fmt.Sprint("User patch failed: ", err))
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
	}

	if uUp.Username != "" {
		// if exists := u.userRepo.ExistByUname(ctx, uUp.Username); exists {
		// 	err := errors.New(fmt.Sprint("Username already taken: ", uname))
		// 	rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		// 	return
		// }
		err := u.userRepo.ChgUsername(ctx, uname, uUp.Username)
		if err != nil {
			u.log.Err("IN [Update] failed to change user username ->", err)
			err = errors.New(fmt.Sprint("User patch failed: ", err))
			rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
			return
		}
	}

	return
}

func (u *userUsecase) ChgPasswd(c context.Context, uuid string, data dtos.ChgPasswdDto) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if data.NPasswd != data.ReNPasswd {
		err := errors.New("New passwords dont match")
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	if exists := u.userRepo.ExistByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("User not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	user, err := u.userRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [ChgPasswd] failed to get user ->", err)
		err = errors.New(fmt.Sprint("User patch failed. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	if auth := u.userRepo.Auth(ctx, user.Username, data.Passwd); !auth {
		err = errors.New("Wrong password")
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	err = u.userRepo.ChgPasswd(ctx, user.Username, data.NPasswd)
	if err != nil {
		u.log.Err("IN [ChgPasswd] failed to change password ->", err)
		err = errors.New(fmt.Sprint("User patch failed: ", err))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *userUsecase) Login(c context.Context, uname string, passwd string) (res domain.User, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if auth := u.userRepo.Auth(ctx, uname, passwd); !auth {
		err := errors.New("Wrong username or password")
		rErr = domain.NewUCaseErr(http.StatusUnauthorized, err)
		return
	}

	user, err := u.userRepo.GetByUsername(ctx, uname)
	if err != nil {
		u.log.Err("IN [Login] failed to get user -> ", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = user

	return
}
