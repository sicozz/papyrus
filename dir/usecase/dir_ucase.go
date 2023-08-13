package usecase

import (
	"context"
	"net/http"
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
		u.log.Err("IN [GetAll]: could not get users ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	return
}
