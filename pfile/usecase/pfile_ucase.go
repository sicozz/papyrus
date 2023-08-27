package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/domain/mapper"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type pFileUseCase struct {
	pFileRepo      domain.PFileRepository
	contextTimeout time.Duration
	log            utils.AggregatedLogger
}

// NewPFileUsecase will create a new dirUsecase object representation of domain.PFileUsecase interface
func NewPFileUsecase(pfr domain.PFileRepository, timeout time.Duration) domain.PFileUsecase {
	logger := utils.NewAggregatedLogger(constants.Usecase, constants.PFile)
	return &pFileUseCase{
		pFileRepo:      pfr,
		contextTimeout: timeout,
		log:            logger,
	}
}

func (u *pFileUseCase) GetAll(c context.Context) (res []dtos.PFileGetDto, rErr domain.RequestErr) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	pFiles, err := u.pFileRepo.GetAll(ctx)
	if err != nil {
		u.log.Err("IN [GetAll] failed to get dirs ->", err)
		rErr = domain.NewUCaseErr(http.StatusInternalServerError, err)
		return
	}

	pFilesDtos := make([]dtos.PFileGetDto, len(pFiles), len(pFiles))
	for i, pf := range pFiles {
		pFilesDtos[i] = mapper.MapPFileToPFileGetDto(pf)
	}

	res = pFilesDtos

	return
}
