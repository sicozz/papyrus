package usecase

import (
	"context"
	"time"

	"github.com/sicozz/papyrus/domain"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	roleRepo       domain.RoleRepository
	contextTimeout time.Duration
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (u *userUsecase) fillUserDetails(ctx context.Context, users []domain.User) (res []domain.User, err error) {
	roles, err := u.roleRepo.GetAll(ctx)
	if err != nil {
		domain.AgLog.Error("Could not get roles to fill user details", err)
	}

	mapRoles := map[int64]domain.Role{}
	for _, role := range roles { //nolint
		mapRoles[role.Code] = role
	}

	// merge the user's data
	for idx, user := range users { //nolint
		if r, ok := mapRoles[user.Role.Code]; ok {
			users[idx].Role = r
		}
	}

	return
}

// NewUserUsecase will create a new userUsecase object representation of domain.UserUsecase interface
func NewUserUsecase(ur domain.UserRepository, rr domain.RoleRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       ur,
		roleRepo:       rr,
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

	// TODO: Check filling details after being able to save users
	domain.AgLog.Info("Not filled:\t", res)
	res, err = u.fillUserDetails(ctx, res)
	domain.AgLog.Info("Filled:\t", res)
	return
}
