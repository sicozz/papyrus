package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/utils"
	"github.com/sicozz/papyrus/utils/constants"
)

type PFileHandler struct {
	PFUsecase domain.PFileUsecase
	log       utils.AggregatedLogger
}

func NewPFileHandler(e *echo.Echo, uu domain.PFileUsecase) {
	logger := utils.NewAggregatedLogger(constants.Delivery, constants.PFile)
	handler := &PFileHandler{uu, logger}
	e.GET("/file", handler.GetAll)
}

func (h *PFileHandler) GetAll(c echo.Context) error {
	h.log.Inf("REQ: get all")
	ctx := c.Request().Context()
	users, rErr := h.PFUsecase.GetAll(ctx)
	if rErr != nil {
		errBody := dtos.NewErrDto("User fetch failed")
		return c.JSON(rErr.GetStatus(), errBody)
	}

	return c.JSON(http.StatusOK, users)
}
