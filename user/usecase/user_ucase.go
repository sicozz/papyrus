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

func (u *userUsecase) GetAll(c context.Context) (res []dtos.UserGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get users ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = make([]dtos.UserGetDto, len(users), len(users))
	for i, u := range users {
		res[i] = mapper.MapUserToUserGetDto(u)
	}
	return
}

func (u *userUsecase) GetByUsername(c context.Context, uname string) (res dtos.UserGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistsByUname(ctx, uname); !exists {
		err := errors.New(fmt.Sprint("User not found. username: ", uname))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	user, err := u.userRepo.GetByUsername(ctx, uname)
	if err != nil {
		u.log.Err("IN [GetByUsername] failed to get user ->", err)
		err = errors.New(fmt.Sprint("User fetch failed. username: ", uname))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return dtos.UserGetDto{}, rErr
	}

	res = mapper.MapUserToUserGetDto(user)

	return
}

func (u *userUsecase) Store(c context.Context, p dtos.UserStore) (res dtos.UserGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user := mapper.MapUserStoreDtoToUser(p)

	if exists := u.userRepo.ExistsByUname(ctx, user.Username); exists {
		err := errors.New("Username already taken")
		rErr = domain.NewUCaseErr(http.StatusConflict, err)
		return
	}

	if exists := u.userRepo.ExistsByEmail(ctx, user.Email); exists {
		err := errors.New("Email already taken")
		rErr = domain.NewUCaseErr(http.StatusConflict, err)
		return
	}

	roleDesc := defRoleDesc
	if p.Role != "" {
		roleDesc = p.Role
	}
	r, err := u.roleRepo.GetByDescription(ctx, roleDesc)
	if err != nil {
		err = errors.New("Role not found")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}
	user.Role = r


	stateDesc := defUserStateDesc
	if p.State != "" {
		stateDesc = p.State
	}
	s, err := u.userStateRepo.GetByDescription(ctx, stateDesc)
	if err != nil {
		err = errors.New("User_state not found")
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}
	user.State = s

	err = u.userRepo.Store(ctx, &user)
	if err != nil {
		err = errors.New("User creation failed")
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	res = mapper.MapUserToUserGetDto(user)

	return
}

func (u *userUsecase) Delete(c context.Context, uname string) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistsByUname(ctx, uname); !exists {
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

func (u *userUsecase) Update(c context.Context, uuid string, p dtos.UserUpdateDto) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if exists := u.userRepo.ExistsByUuid(ctx, uuid); !exists {
		err := errors.New(fmt.Sprint("User not found. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusNotFound, err)
		return
	}

	user, err := u.userRepo.GetByUuid(ctx, uuid)
	if err != nil {
		u.log.Err("IN [Update] failed to get user ->", err)
		err = errors.New(fmt.Sprint("User patch failed. uuid: ", uuid))
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	taken := u.userRepo.ExistsByUname(ctx, p.Username)
	if taken && p.Username != user.Username {
		err := errors.New(fmt.Sprint("Username already taken: ", p.Username))
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	taken = u.userRepo.ExistsByEmail(ctx, p.Email)
	if taken && p.Email != user.Email {
		err := errors.New(fmt.Sprint("Email already taken: ", p.Email))
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	if p.Role != "" {
		if exists := u.roleRepo.ExistsByDescription(ctx, p.Role); !exists {
			err := errors.New(fmt.Sprint("Role not found. description: ", p.Role))
			rErr = domain.NewUCaseErr(http.StatusNotFound, err)
			return
		}
	}

	if p.State != "" {
		if exists := u.userStateRepo.ExistsByDescription(ctx, p.State); !exists {
			err := errors.New(fmt.Sprint("State not found. description: ", p.State))
			rErr = domain.NewUCaseErr(http.StatusNotFound, err)
			return
		}
	}

	err = u.userRepo.Update(ctx, uuid, p)
	if err != nil {
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}

func (u *userUsecase) ChgPasswd(c context.Context, uuid string, data dtos.UserChgPasswdDto) (rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if data.NPasswd != data.ReNPasswd {
		err := errors.New("New passwords dont match")
		rErr = domain.NewUCaseErr(http.StatusNotAcceptable, err)
		return
	}

	if exists := u.userRepo.ExistsByUuid(ctx, uuid); !exists {
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

func (u *userUsecase) Login(c context.Context, uname string, passwd string) (res dtos.UserGetDto, rErr domain.RequestErr) {
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

	res = mapper.MapUserToUserGetDto(user)

	return
}
